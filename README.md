# FileVer
-------------------------------------

FileVer (File Version) automatically versions files for packaging and
distribution. File versioning is essential for browser cache busting for
Javascript, HTML, CSS, images, and other assets loaded by the browser.  

FileVer has two main functions:
 1. (Version) Hash versioned files and generate FileVer. Place new versioned
    files into an output directory. Delete any old versions of that file.  
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
The default version is 32 bits, 8 character version with a 1 out of 4.2 billion
chance for collision.   (We considered 7 characters, 1 out of 268 million
chance, but it's better to use a whole byte instead of 5.25 bytes. We consider 6
characters to be too small with a 1 out of 16 million chance of collision.)

The whole version is 12 characters long; the deliminator `~` is 1 character ,
identifier `fv` is 2 characters, key/value delimiter `=` is 1 character, and 8
characters for the base64 version.

	app~fv=4mIbJJPq.min.js


# Naming
A FileVer is the full file name with a file version.

General naming (not specific to FileVer)
| Name                                  | Example                   |
| ------------------------------------- | ------------------------- |
| Full path                             | e/app.min.js              |
| Directory (dir)                       | e/                        |
| File (filename)                       | app.min.js                |
| Base                                  | app                       |
| Extension                             | .min.js                   |
| Base extension                        | .js                       |

Naming for a FileVer example
| Name                                  | Example                   |
| ------------------------------------- | ------------------------- |
| Full path                             | e/app~fv=4mIbJJPq.min.js  |
| Dir                                   | e/                        |
| File (filename)                       | app~fv=4mIbJJPq.min.js    |
| Base                                  | app~fv=4mIbJJPq           |
| Extension                             | .min.js                   |
| Base extension                        | .js                       |
|                                       |                           |
| Pathed FileVer                        | e/app~fv=4mIbJJPq.min.js  |
| FileVer, versioned file               | app~fv=4mIbJJPq.min.js    |
| Bare path                             | e/app.min.js              |
| Bare file                             | app.min.js                |
| Bare                                  | app                       |
| Version                               | 4mIbJJPq                  |
| Delimiter (delim)                     | ~fv=                      |
| Delimited (delim'd) version           | ~fv=4mIbJJPq              |
| Dummy version (zeroed version)        | 00000000                  |
| Delim'd dummy Version                 | ~fv=00000000              |
| Dummy versioned file (full dummy)     | e/app~fv=00000000.min.js  |



	URI Specific
		https://example.com:8081/bob/joe.txt?name=ferret#nose?name=bob
		\___/   \______________/\__________/ \_________/ \___________/
		 |            |              |            |           |
		scheme     authority        path        query      fragment
		        \_________/\___/                         \__/\_______/
		           |       |                              |      |
		         host     port                          anchor  fquery
		                                     \_______________________/
		                                                |
		                                              quag


Naming for URI Paths
| Name                 | Example                     |
| -------------------- | --------------------------- |
| Full path            | https://cyphr.me:8081/assets/img/cyphrme_long.png  |
| Scheme               | https:                      |
| Authority            | cyphr.me:8081               |
| Host                 | cyphr.me                    |
| Port                 | :8081                       |
| URIPath              | bob/joe.txt                 |
| Query                | name=ferret                 |
| Fragment             | nose?name=bob               |
| Anchor               | nose                        |
| FragmentQuery        | ?name=bob                   |
| Quag                 | ?name=ferret#nose?name=bob  |

Additionally, the normal path information will be populated. 
| Name                 | Example                     |
| -------------------- | --------------------------- |
| Directory (dir)      | https://cyphr.me:8081/assets/img/    |
| File (filename)      | cyphrme_long.png                     |
| Base                 | cyphrme_long                         |
| Extension            | .png                                 |
| Base extension       | .png                                 |



A FileVer may refer to a file with a dummy version, i.e. `app~fv=00000000.min.js`
is a FileVer.  

A base name of a FileVer is everything before the `~` character and excludes
pathing.  For the path `/example/app~fv=4mIbJJPq.min.js`, the base is simply
`app` with no extension or pathing.  

Extension includes all sub extensions and a preceding `.`.  For
`app~fv=4mIbJJPq.min.js.map` the extension is `.min.js.map`.

The base extension is always the last extension if present.  It may be equal to
extension.  For example, a file `app.js` will have a base extension and
extension of `.js`

### Vocabulary
- `dist` -"distribution": The destination directory.  
- Versioned file: A file with a file version, e.g. `app~fv=4mIbJJPq.min.js`.
- Delimiter:  The FileVer delimiter string is (by default) `~fv=`.  The ending
  delimiter for a Version is any non-base64 character, such as another
  "~" character.  This follows the standard URL Query and URL Fragment Query
  notation.  
- Non-dummy versioned files: a file with a digest as the version, e.g.
  `app~fv=4mIbJJPq.min.js`.


# Pipeline
## Directory Structure

```
parent_directory
  ⎿ src
  ⎿ dist
```
Input `src` directory structure is preserved in `dist`. FileVer will version
files in `src`, output to `dist`, and perform Replace() in dist. FileVer by
default recreates the directory structure of `src` into `dist`. If this is not
desired, see examples for "manual" versioning.  

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
  - FileVer versions source file and outputs `dist/file~fv=4mIbJJPq.min.js` -> 
  - FileVer updates other source code files in `dist` with FileVer (Replace).


# Dummies - Import References to Versioned Files
All text based source files that refer to versioned file should use the **dummy
version** in import references in the `src` directory.  After running, FileVer
will update references in `dist` with correct versioning (input will be left
untouched).

It is correct for Javascript files in `src` to import use this form: 

`import * as test1 from './test_1~fv=00000000.js';`

FileVer Replace() will update imports in `dist` to point to the correct file.

`import * as test1 from './test_1~fv=4WYoW0MN.js';`


| `src` (Input Directory) Import   | `dist` (Output Directory) Import |
| -------------------------------- | -------------------------        |
| `app~fv=00000000.min.js`         | `app~fv=4mIbJJPq.min.js`         |
| `lib~fv=00000000.min.js`         | `lib~fv=820OsC4y.min.js`         | 


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
`app~fv=00000000.min.js`.  Then, after the file version digest is calculated and
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
