# C/C++ Lockfile-Less Scanning

In many C/C++ projects, libraries are often included directly in the source code folder. Our product offers a unique lockfile-less scanning method to address this issue effectively.

### Benefits <a href="#benefits" id="benefits"></a>

* **Comprehensive Coverage**: Unlike most SAST tools, Aikido scans C/C++ libraries included in the source code folder. This ensures that no dependencies go untracked, enhancing the security of your codebase.
* **Accurate Matching**: By hashing all relevant files in the repository, we can accurately match your project's dependencies with a database of the most popular open-source libraries.
* **Detection of Modified Libraries**: Even if you have made small edits to a library or if a small file is missing, our tool can still detect the library version and match it correctly.

### How It Works <a href="#how-it-works" id="how-it-works"></a>

1. **Database of Open Source Libraries**: We maintain a large database of popular open-source libraries, including all its versions and files.
2. **Hashing Source Files**: When a repository is scanned, our tool hashes all relevant C/C++ files in the project. These hashes are then compared with the hashes in our database.
3. **Matching and Detection**: If there are sufficient matches between the hashes from the project and those in the database, we identify the library versions that are used.
