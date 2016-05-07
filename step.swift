#!/usr/bin/swift

import Foundation

typealias ArgsArray = Array<String>

let carthageDirName = "Carthage"
let buildDirName = "Build"
let cacheFileName = "Cachefile"
let resolvedFileName = "Cartfile.resolved"

let env = NSProcessInfo.processInfo().environment
let task = NSTask()

guard let workingDir = env["working_dir"] where workingDir != "" else {
    print("Working directory set to empty string, or nil. Exiting.. Please set a working directory and re-run the build.")
    exit(0)
}

guard let carthageCommand = env["carthage_command"] else {
    fatalError("no command to execute")
}

let bootstrapCommand = carthageCommand == "bootstrap"

func collectArgs(env: [String : String]) -> ArgsArray {
    var args = ArgsArray()
    
    if let platform = env["platform"] {
        args.append("--platform " + platform)
    }
    
    if let verboseOutput = env["verbose_output"] where verboseOutput == "true" {
        args.append("--verbose")
    }
    
    if let noUseBinaries = env["no_use_binaries"] where noUseBinaries == "true" {
        args.append("--no-use-binaries")
    }
    
    if let sshOutput = env["ssh_output"] where sshOutput == "true" {
        args.append("--use-ssh")
    }
    
    if let carthageOptions = env["carthage_options"] {
        args.append(carthageOptions)
    }
    
    return args
}

func swiftVersion() -> String? {
    let swiftVersionTask = NSTask()
    swiftVersionTask.currentDirectoryPath = workingDir
    swiftVersionTask.launchPath = "/usr/bin/swift"
    swiftVersionTask.arguments = ["-version"]
    
    let pipe = NSPipe()
    swiftVersionTask.standardOutput = pipe
    swiftVersionTask.launch()
    
    let data = pipe.fileHandleForReading.readDataToEndOfFile()
    guard let versionString = String(data: data, encoding: NSUTF8StringEncoding) else {
        return nil
    }

    return versionString
}

func contentsOfCartfileResolved() -> String? {
    guard let cartfileResolvedData = NSFileManager.defaultManager().contentsAtPath("\(workingDir)/\(resolvedFileName)"),
        let cartfileResolvedContent = String(data: cartfileResolvedData, encoding: NSUTF8StringEncoding) else {
            return nil
    }

    return cartfileResolvedContent
}

func cacheAvailable() -> Bool {
    if !NSFileManager.defaultManager().fileExistsAtPath("\(workingDir)/\(carthageDirName)") {
        return false
    }

    do {
        try NSFileManager.defaultManager().contentsOfDirectoryAtPath("\(workingDir)/\(carthageDirName)/\(buildDirName)")
    } catch _ {
        return false
    }

    // read cache
    guard let cacheFileData = NSFileManager.defaultManager().contentsAtPath("\(workingDir)/\(carthageDirName)/\(cacheFileName)"),
        let cacheFileContents = String(data: cacheFileData, encoding: NSUTF8StringEncoding),
        let version = swiftVersion(),
        let resolved = contentsOfCartfileResolved() else {
            return false
    }
    
    let contents = "--Swift version: \(version) --Swift version \n --\(resolvedFileName): \(resolved) --\(resolvedFileName)"

    return cacheFileContents == contents
}

// exit if bootstrap is cached
let hasCachedItems = cacheAvailable()
if bootstrapCommand && hasCachedItems {
    print("Cache available for bootstrap command, exiting. If you would like to update your Carthage contents, select `update` as Carthage command and re-run your build.")
    exit(0)
}

let command = "carthage \(carthageCommand)"
var args = " " + ( collectArgs(env).map { "\($0)" } ).joinWithSeparator(" ")

task.currentDirectoryPath = workingDir
task.launchPath = "/bin/bash"
task.arguments = ["-c", command + args]

print("Running carthage command: \(task.arguments!.reduce("") { str, arg in str + "\(arg) " })")

task.launch()
task.waitUntilExit()

// create cache
if bootstrapCommand {
    let cacheFilePath = "\(workingDir)/\(carthageDirName)/\(cacheFileName)"
    guard let version = swiftVersion(),
        let resolved = contentsOfCartfileResolved() else {
            print("Failed to create cache content.")
            exit(0)
    }
    
    let contents = "--Swift version: \(version) --Swift version \n --\(resolvedFileName): \(resolved) --\(resolvedFileName)"

    if NSFileManager.defaultManager().fileExistsAtPath("\(workingDir)/\(carthageDirName)") {
        do {
            try contents.writeToFile(cacheFilePath, atomically: false, encoding: NSUTF8StringEncoding)
        } catch _ {
            print("Failed to update CacheFile.")
            exit(0)
        }
    } else {
        // create Cachefile
        if NSFileManager.defaultManager().createFileAtPath(cacheFilePath, contents: contents.dataUsingEncoding(NSUTF8StringEncoding), attributes: nil) {
            print("Cachefile created successfully.")
        } else {
            print("Failed to create Cachefile.")
        }
    }
}
