#!/usr/bin/swift

import Foundation

typealias ArgsArray = Array<String>

let carthageDirName = "Carthage"
let buildDirName = "Build"
let cacheFileName = "Cachefile"
let resolvedFileName = "Cartfile.resolved"

let env = ProcessInfo.processInfo.environment
let task = Process()

guard let workingDir = env["working_dir"], workingDir != "" else {
    print("Working directory set to empty string, or nil. Exiting.. Please set a working directory and re-run the build.")
    exit(1)
}

guard let carthageCommand = env["carthage_command"] else {
    fatalError("no command to execute")
}

let bootstrapCommand = carthageCommand == "bootstrap"
let checkoutCommand = carthageCommand == "checkout"

func collectArgs(_ env: [String : String]) -> ArgsArray {
    var args = ArgsArray()

    if let platform = env["platform"] {
        args.append("--platform " + platform)
    }

    if let verboseOutput = env["verbose_output"], verboseOutput == "true" {
        args.append("--verbose")
    }

    if let noUseBinaries = env["no_use_binaries"], noUseBinaries == "true" {
        args.append("--no-use-binaries")
    }

    if let sshOutput = env["ssh_output"], sshOutput == "true" {
        args.append("--use-ssh")
    }

    if let carthageOptions = env["carthage_options"] {
        args.append(carthageOptions)
    }

    return args
}

func swiftVersion() -> String? {
    let swiftVersionTask = Process()
    swiftVersionTask.currentDirectoryPath = workingDir
    swiftVersionTask.launchPath = "/usr/bin/swift"
    swiftVersionTask.arguments = ["-version"]

    let pipe = Pipe()
    swiftVersionTask.standardOutput = pipe
    swiftVersionTask.launch()

    let data = pipe.fileHandleForReading.readDataToEndOfFile()
    guard let versionString = String(data: data, encoding: String.Encoding.utf8) else {
        return nil
    }

    return versionString
}

func contentsOfCartfileResolved() -> String? {
    guard let cartfileResolvedData = FileManager.default.contents(atPath: "\(workingDir)/\(resolvedFileName)"),
        let cartfileResolvedContent = String(data: cartfileResolvedData, encoding: String.Encoding.utf8) else {
            return nil
    }

    return cartfileResolvedContent
}

func cacheAvailable() -> Bool {
    if !FileManager.default.fileExists(atPath: "\(workingDir)/\(carthageDirName)") {
        return false
    }

    do {
        try FileManager.default.contentsOfDirectory(atPath: "\(workingDir)/\(carthageDirName)/\(buildDirName)")
    } catch _ {
        return false
    }

    // read cache
    guard let cacheFileData = FileManager.default.contents(atPath: "\(workingDir)/\(carthageDirName)/\(cacheFileName)"),
        let cacheFileContents = String(data: cacheFileData, encoding: String.Encoding.utf8),
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
    print("Using cached dependencies for bootstrap command. If you would like to force update your dependencies, select `update` as Carthage command and re-run your build.")
    exit(0)
}

let command = "carthage \(carthageCommand)"
var args = " "

if !checkoutCommand {
    args = args + ( collectArgs(env).map { "\($0)" } ).joined(separator: " ")
}


task.currentDirectoryPath = workingDir
task.launchPath = "/bin/bash"
task.arguments = ["-c", command + args]

print("Running carthage command: \(task.arguments!.reduce("") { str, arg in str + "\(arg) " })")

task.launch()
task.waitUntilExit()

guard task.terminationStatus == 0 else {
    exit(task.terminationStatus)
}

// create cache
if bootstrapCommand {
    let cacheFilePath = "\(workingDir)/\(carthageDirName)/\(cacheFileName)"
    guard let version = swiftVersion(),
        let resolved = contentsOfCartfileResolved() else {
            print("Failed to create cache content.")
            exit(1)
    }

    let contents = "--Swift version: \(version) --Swift version \n --\(resolvedFileName): \(resolved) --\(resolvedFileName)"

    if FileManager.default.fileExists(atPath: "\(workingDir)/\(carthageDirName)") {
        do {
            try contents.write(toFile: cacheFilePath, atomically: false, encoding: String.Encoding.utf8)
        } catch _ {
            print("Failed to update CacheFile.")
            exit(1)
        }
    } else {
        // create Cachefile
        if FileManager.default.createFile(atPath: cacheFilePath, contents: contents.data(using: String.Encoding.utf8), attributes: nil) {
            print("Cachefile created successfully.")
        } else {
            print("Failed to create Cachefile.")
            exit(1)
        }
    }
}

exit(0)
