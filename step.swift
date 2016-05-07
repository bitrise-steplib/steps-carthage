#!/usr/bin/swift

import Foundation

typealias ArgsArray = Array<String>

let carthageDirName = "Carthage"
let buildDirName = "Build"
let cacheFileName = "CacheFile"

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

let env = NSProcessInfo.processInfo().environment
let task = NSTask()

var cacheAvailable = false

if let workingDir = env["working_dir"] where workingDir != "" {
    task.currentDirectoryPath = workingDir

    func fileExistsAtCarthagePath(pathComponent: String) -> Bool {
        return NSFileManager.defaultManager().fileExistsAtPath("\(workingDir)/\(pathComponent)")
    }

    let carthageDir = fileExistsAtCarthagePath(carthageDirName)
    if carthageDir {
        let buildDir = fileExistsAtCarthagePath("\(carthageDirName)/\(buildDirName)")
        let buildDirContents: [NSString]?

        do {
            buildDirContents = try NSFileManager.defaultManager().contentsOfDirectoryAtPath("\(workingDir)/\(carthageDirName)/\(buildDirName)")
        } catch _ {
            buildDirContents = nil
        }
        if buildDir && buildDirContents != nil {
            cacheAvailable = true
        }
    }
}

guard let carthageCommand = env["carthage_command"] else {
    fatalError("no command to execute")
}

// read cache
if carthageCommand == "bootstrap" && cacheAvailable {
    print("Cache available for bootstrap command, exiting. If you would like to update your Carthage contents, select `update` as Carthage command and re-run your build.")
    exit(0)
}

let command = "carthage \(carthageCommand)"
var args = " " + ( collectArgs(env).map { "\($0)" } ).joinWithSeparator(" ")

task.launchPath = "/bin/bash"
task.arguments = ["-c", command + args]

print("Running carthage command: \(task.arguments!.reduce("") { str, arg in str + "\(arg) " })")

// run the shell command
task.launch()
//
// // ensure to be finished before another command can run
task.waitUntilExit()
