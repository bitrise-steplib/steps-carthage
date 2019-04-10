# Carthage step for iOS/Mac

Plain simple build step for setting up your project with [Carthage](https://github.com/Carthage/Carthage).

## Cache Carthage dependencies on [Bitrise](https://www.bitrise.io)

You can cache the result of `carthage bootstrap` to make Bitrise faster when building unchanged dependencies. For that you'll have to [set up caching](http://devcenter.bitrise.io/caching/about-caching/) in your Bitrise workflow.

1. add a `Cache pull` step **before** your `Carthage` step,
2. add a `Cache push` step **after** your `Carthage` step,

The `Cachefile` stores a Swift version you ran `carthage bootstrap` the last time and the content of your `Cartfile.resolved`. Until either of these information is not changed between builds, Bitrise will ignore the `bootstrap` call and use the cached content of your `Carthage/Build` directory for building your project. If you have changes in your `Cartfile.resolved`, or changed the stack to one with a different Swift version, it will run `carthage bootstrap` to make sure the cache is only used when it's 100% compatible.

## Run locally

Can be run directly with the [bitrise CLI](https://github.com/bitrise-io/bitrise),
just `git clone` this repository, `cd` into it's folder in your Terminal/Command Line
and call `bitrise run test`.

*Check the `bitrise.yml` file for required inputs which have to be
added to your `.bitrise.secrets.yml` file!*
