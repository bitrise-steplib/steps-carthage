## Changelog (Current version: 3.0.5)

-----------------

### 3.0.5 (2017 Feb 01)

* [f69851c] prepare for 3.0.5
* [cbf38c9] Work dir (#36)

### 3.0.4 (2017 Jan 11)

* [f0f15f3] prepare for 3.0.4
* [e74a19c] default command changed to "bootstrap" (#33)
* [6d13955] Cache setup link fix

### 3.0.3 (2017 Jan 03)

* [5c17351] prepare for 3.0.3
* [d3384a3] Update main.go (#32)

### 3.0.2 (2016 Dec 19)

* [8e1657d] prepare for 3.0.2
* [f0a09fd] add macos tag (#31)

### 3.0.1 (2016 Dec 12)

* [b3720bb] prepare for 3.0.1
* [910c7f9] Fix typo (#29)

### 3.0.0 (2016 Nov 18)

* [1c7f782] prepare for 3.0.0
* [7298a9b] Go toolkit (#28)

### 2.4.0 (2016 Sep 26)

* [4961160] prepare for 2.4.0
* [2777a0d] GITHUB_ACCESS_TOKEN support (#25)

### 2.3.0 (2016 Sep 12)

* [32e45fc] prep for v2.3.0
* [825d6bd] Use step-swift3.swift when Swift3 (#24)

### 2.2.0 (2016 Aug 09)

* [6fc5204] removed unused bump version workflow
* [7efd34f] prep for v2.2.0
* [456ed16] Merge pull request #23 from KoCMoHaBTa/master
* [3983e0e] addressed PR comments - removed unnecessary empty line
* [9b80d83] added support for carthage checkout command

### 2.1.2 (2016 Jul 12)

* [aee1c8a] prepare for 2.1.2
* [17c987c] Merge pull request #22 from bitrise-steplib/review
* [146b898] review

### 2.1.1 (2016 May 18)

* [21df866] Merge pull request #21 from dkalachov/exitstatus
* [567e1f4] Do not exit with 0 code on error.

### 2.1.0 (2016 May 07)

* [6c21190] Merge pull request #20 from vasarhelyia/feature/cache
* [4744b50] Fix typo
* [30835f1] Update readme with cache and cache related log messages
* [aaa1311] Reorganize content reads to eliminate SIL segfaults
* [76ad2dc] Add const for resolved file name
* [cbade3d] Fix typos
* [08dd1dc] Logic for caching
* [5a0a73c] Check for existing cache

### 2.0.2 (2016 Apr 25)

* [9b01839] Merge pull request #19 from akashivskyy/master
* [4483abd] Add default value for work dir
* [b97d90c] Fix empty working dir bug

### 2.0.1 (2016 Mar 01)

* [703f0c9] Use string concat for args
* [6569b92] Printing carthage command to run

### 2.0.0 (2016 Feb 16)

* [11a34ef] Merge pull request #18 from vasarhelyia/swiftification
* [544225a] Pass args as array
* [c84adc4] Use guard
* [9b0ba29] Eliminate duplicated dashes
* [c28d64e] Comment out input vars
* [dec34b1] set -e
* [50f6ec9] Indentation fixes
* [3839bce] Add swift step script

### 1.0.8 (2016 Feb 04)

* [1effd0f] Fix yml syntax
* [7b0f353] Fix use new deps syntax
* [6bae536] Merge pull request #17 from vasarhelyia/use-preinstalled-carthage
* [2eca4da] Remove explicit install
* [8ad166b] Add carthage as dependency

### 1.0.7 (2016 Jan 27)

* [833b276] Merge pull request #16 from toshi0383/ts-add-platforms
* [efbce4f] Add tvOs and all to --platform option

### 1.0.6 (2016 Jan 09)

* [c116e56] Make update the default step
* [abeb6c7] Add CLI usage explanation
* [ee28b8e] Add basic test workflow
* [c3b64f0] Add gitignore
* [0715a38] Add bitrise.yml
* [e758c86] Add general options

### 1.0.5 (2015 Dec 13)

* [6e5c311] Fix indentation
* [3efec75] Merge pull request #14 from julio-rivera/master
* [51cbbac] Add ssh mode

### 1.0.4 (2015 Nov 28)

* [f955d6f] Merge pull request #12 from vasarhelyia/feature/platform-as-input
* [c926437] Add missing param name
* [84fd462] Set ex output
* [2af5a19] Remove extraneous space from command param list
* [01ff8f0] Remove extra string from command
* [278337e] Add platform as input var

### 1.0.3 (2015 Nov 25)

* [5825a8c] Merge pull request #11 from vasarhelyia/fix/remove-outdated-homebrew-formula-workaround
* [159a77d] Remove extra input for carthage version

### 1.0.2 (2015 Nov 06)

* [daf9135] Merge pull request #10 from viktorbenei/patch-1
* [b248f1e] Fix for workdir switch

### 1.0.1 (2015 Oct 27)

* [4011cb9] Update README.md
* [d78d9a5] Merge pull request #8 from vasarhelyia/feature/specify-working-directory
* [f32590c] Add input for working directory
* [4d686af] Merge pull request #7 from vasarhelyia/feature/specify-carthage-version
* [f19e7ef] Add input to use brew default carthage or install 0.9.4 from package
* [7c67b8c] Merge pull request #6 from vasarhelyia/feature/verbose-mode
* [2d58e40] Add verbose run input
* [a390264] Set optional verbose mode to carthage command

-----------------

Updated: 2017 Feb 01