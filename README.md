# FileVer
-------------------------------------

FileVer (File Version) automatically versions files for packaging and
distribution. File versioning is essential for browser cache busting for
Javascript, HTML, CSS, images, and other assets loaded by the browser.  

FileVer has two main functions:
 1. (Version) Hash versioned files to generate the file version. Placing
    versioned files, named with the file version, into an output directory.
    Delete any old versions of that file.  
 2. (Replace) In the output directory (`dist`), update references in source
    files to versioned files.

We recommend using FileVer in conjunction with [watch][watch] and once
configured, a file change will trigger a global update of file versions and
references in source files, making development fast and painless. Compared to
other methods, FileVer shines when source files refer to versioned files and
when file contents are frequently updated and changed.

FileVer may be used in pipelines that involves these keywords:
TypeScript, Sass, Go Templates, handlebars, Javascript minification, HTML
minification, and CSS minification.


# Version
The version string is a digest derived from the content of the file.  
The default version is 48 bits, 12 character long (8 characters of base64 URI),
with a 1 out of 281 trillion chance for collision. (We considered 7 characters,
1 out of 4 trillion chance, but that's 5.25 bytes and it's better to use a whole
byte.)

	app~fv=4mIbJJPq.min.js

A base file of a FileVer is everything before the last `~` character. 

# Naming
A FileVer is the full file name with a file version.

| Name                                  | Example                   |
| ------------------------------------- | ------------------------- |
| FileVer, versioned file, full version | `app~fv=4mIbJJPq.min.js`  |
| File, File name, bare file            | `app.min.js`              |
| Version                               | `4mIbJJPq`                |
| Delimiter, delim                      | `~fv=`                    |
| Delim'd Version                       | `~fv=4mIbJJPq`            |
| zero'd Version, dummy Version         | `00000000`                |
| Delim'd dummy Version                 | `~fv=00000000`            |
| Dummy Versioned File, Full Dummy      | `app~fv=00000000.min.js`  |

A FileVer may refer to a file with a dummy version, i.e. app~fv=00000000.min.js is a FileVer.  

### Vocabulary
- dist -"distribution": The destination directory.  
- Versioned file: A file with a file version, e.g. app~fv=4mIbJJPq.min.js.
- Delimitor:  The FileVer delimiter string is (by default) `~fv=`.  The ending
  delimitor for a Version is any non-base64 character, such as another
  "~" character.  This follows the standard URL Query and URL Fragment Query
  notation.  
- non-dummy versioned files: a file with a digest as the version, e.g.
  app.min.js~fv=4mIbJJPq.


# Pipeline
## Directory Structure

```
your_directory
  ??? src
  ??? dist
```
Input `src` directory structure is preserved in `dist`.
FileVer will version files in `src`, output to `dist`, and perform Replace() in
dist. FileVer by default recreates the directory structure of `src` into `dist`.
If this is not desired, see examples for "manual" versioning.  

Files that are not concerned with file versioning should be placed directly in
`dist`.  

Directories `src` and `dist` may be named as desired.  


## Development Pipeline
The suggested pipeline is to configure `watch` to watch relevant source files.
1. Use [watch][watch] to watch specific files for change.
2. On file change, configure `watch` to run a script that does (esbuild ->
FileVer). FileVer is responsible for hashing the updated file, placing it into
the dist directory, and update any source files references to the updated file
in `dist`.

Then:
  - modify js source file (file.js)-> 
  - watch is triggered, runs a `.sh` script that invokes 1. esbuild and then 2. FileVer ->
  - esbuild minifies source file, outputs `src/file.min.js`. 
  - FileVer versions source file and outputs `dist/file.min.js~fv=4mIbJJPq` -> 
  - FileVer updates other source code files in `dist` with FileVer (Replace).


# Dummies - Import References to Versioned Files
All text based source files that refer to versioned file should use the **dummy
version** in import references in the `src` directory.  After running, FileVer
will update references in `dist` with correct versioning (input will be left
untouched).

It is correct for Javascript files in `src` to import use this form: 

`import * as test1 from './test_1.js~fv=00000000';`

FileVer Replace() will update imports in `dist` to point to the correct file.

`import * as test1 from './test_1.js~fv=4WYoW0MN';`


| `src` (Input Directory) Import   | `dist` (Output Directory) Import |
| -------------------------------- | -------------------------        |
| `app.min.js~fv=00000000`         | `app.min.js~fv=4mIbJJPq`         |
| `lib.min.js~fv=00000000`         | `lib.min.js~fv=820OsC4y`         | 


Note: A consequence of this design is that the versioned files digest will not
match the current file name unless the references to versioned files are
re-zeroed.  FileVer design assumes all Versioned file imports are zero'd.  

## Not using dummy files, not using `src`, or not using zero'd imports.  
If not wanting to use dummy files, each to-be-versioned file must be
individually enumerated.  See `Example_noDummy()` for a demonstration. Replace()
will still be needed to update references in `dist`.

## All source file imports must be relative to root directory 
When contemplating the design, there were two designs considered: 
  1. All files names had to be global unique or 
  2. Use paths for namespacing.  
	
The second option allows "duplicate" file names using path as a namespace, and
is more Unix-link.  There are many advantages to 2 over 1, so design 2 was
implemented.  This also allows the SAVR regex to be global, instead of needing
to generate a SAVR for each subdirectory.  

To implement 2 easily, all imports in `dist` must be relative to directory
`dist`.  

For example, if in subdirectory named "subdir", imports must be relative to the root, e.g. 

'../subdir/test_3~fv=CX_w_yNh.js';

and not

'./test_3~fv=00000000.js'

 
#### Matching problem if not relative to root

Given the files 

```
test_3.js
subdir/test_3.js
```

Although these are easy to match in root, there is not good way to match a
secondary SAVR from `subdir`.  There is no good way to write a regex to match
the correct file.  

```
'../test_3.js';
'./test_3.js';
```

Any regex constrain would put restrictions on the syntax of the language itself.
We want FileVer to be syntax agnostic, other than Unix pathing which is widely
supported.  Including relative root significantly simplifies.  

## Update Recursion
The suggested version pipeline uses an input and output directory in order to
avoid update recursion. Since some formats, like Javascript modules, may refer
to other versioned files, this risks file version recursion. To avoid this, the
pipeline is one directional. It's suggested that source files that refer to
other versioned files use a dummy, constant file version in references, such as
`app.min.js~fv=00000000`.  Then, after the file version digest is calculated and
placed in the dist dir, the version in other source files is updated with
Replace().  This has the consequence that the versioned file's name, after
replace, will not be calculable from the source unless the dummy version is
restored.

# Examples
See [version_test.go](version_test.go)

# FAQ

## Why not use HTTP ETag?
Because [ETags][ETag] require HTTP requests.  Versioning precludes any HTTP
request after the initial page load.  That's far more scalable and efficient,
especially considering the HTTP request itself is one of the most costly parts
of page load.  

# See also:
DirHash: https://github.com/golang/mod/blob/ce943fd02449f621243c9ea6e64098e84752b92b/sumdb/dirhash/hash.go



----------------------------------------------------------------------
# Attribution, Trademark notice, and License
FileVer is released under The 3-Clause BSD License. 

"Cyphr.me" is a trademark of Cypherpunk, LLC. The Cyphr.me logo is all rights
reserved Cypherpunk, LLC and may not be used without permission.


[watch]: https://github.com/Cyphrme/watch
[ETag]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/ETag
