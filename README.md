# FileVer
-------------------------------------

File versioning package `filever` automatically versions files for packaging and
distribution. File versioning is essential for browser cache busting for
Javascript, HTML, CSS, images, and other assets loaded by the browser.  

`filever` has two main functions:
 1. (Version) Hash versioned files and generate filever. Place new versioned
    files into an output directory. Delete any old versions of that file.  
 2. (Replace) In the output directory (`dist`), update references in source
    files to versioned files.

We recommend using `filever` in conjunction with [watchmod][watchmod] and once
configured, a file change will trigger a global update of file versions and
references in source files, making development fast and painless. Compared to
other methods, `filever` shines when source files refer to versioned files and
when file contents are frequently updated and changed.

`filever` may be used in pipelines that involves these keywords: TypeScript,
Sass, Go Templates, handlebars, Javascript minification, HTML minification, and
CSS minification.


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
See package [path][path] for file and URI path naming.  A filever is the full
file name with a file version.



General naming not specific to `filever`:
| Name                                  | Example                   |
| ------------------------------------- | ------------------------- |
| Full path                             | e/app~fv=4mIbJJPq.min.js  |
| Dir                                   | e/                        |
| File (filename)                       | app~fv=4mIbJJPq.min.js    |
| FileBase                              | app~fv=4mIbJJPq           |
| Extension                             | .min.js                   |
| Extension base                        | .js                       |

`filever` specific naming
| Name                                  | Example                   |
| ------------------------------------- | ------------------------- |
| Pathed filever                        | e/app~fv=4mIbJJPq.min.js  |
| Filever                               | app~fv=4mIbJJPq.min.js    |
| Bare path                             | e/app.min.js              | 
| Bare file                             | app.min.js                |
| Bare                                  | app                       |
| Version                               | 4mIbJJPq                  |
| Delimiter (delim)                     | ~fv=                      |
| Delimited (delim'd) version           | ~fv=4mIbJJPq              |
| Dummy version (zeroed version)        | 00000000                  |
| Delim'd dummy Version                 | ~fv=00000000              |
| Dummy versioned file (full dummy)     | e/app~fv=00000000.min.js  |

Bare path is the full path stripped of versioning.  

A filever may refer to a file with a dummy version, i.e.
`app~fv=00000000.min.js` is a legitimate filever.  

The base name of a filever is everything before the `~` character and excludes
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
- Delimiter:  The filever delimiter string is (by default) `~fv=`.  The ending
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
1. Use [watchmod][watchmod] to watch specific files for change.
2. On file change, configure `watch` to run a script that does (esbuild ->
`filever`). `filever` is responsible for hashing the updated file, placing it
into the dist directory, and update any source files references to the updated
file in `dist`.

Then:
  - modify js source file (file.js)
  - watchmod is triggered, runs a `.sh` script that invokes 1. esbuild and then
    2. `filever`
  - esbuild minifies source file, outputs `src/file.min.js`
  - `filever` versions source file and outputs `dist/file~fv=4mIbJJPq.min.js`
    (Version)
  - `filever` updates other source code files in `dist` with filever (Replace).


# Dummies - Import References to Versioned Files
All text based source files that refer to versioned file should use the **dummy
version** in import references in the `src` directory.  After running, `filever`
updates references in `dist` with versioning (input will be left untouched).

It is correct for Javascript files in `src` to import use this form: 

`import * as test1 from './test_1~fv=00000000.js';`

`filever` Replace() updates imports in `dist` to point to the correct file.

`import * as test1 from './test_1~fv=4WYoW0MN.js';`


| `src` (Input Directory) Import   | `dist` (Output Directory) Import |
| -------------------------------- | -------------------------        |
| `app~fv=00000000.min.js`         | `app~fv=4mIbJJPq.min.js`         |
| `lib~fv=00000000.min.js`         | `lib~fv=820OsC4y.min.js`         | 


Note: A consequence of this design is that the versioned files digest will not
match the current file name unless the references to versioned files are
re-zeroed.  `filever`'s design assumes all Versioned file imports are zero'd.  

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
We want `filever` to be syntax agnostic, other than Unix pathing which is widely
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
`filever` is released under The 3-Clause BSD License. 

"Cyphr.me" is a trademark of Cypherpunk, LLC. The Cyphr.me logo is all rights
reserved Cypherpunk, LLC and may not be used without permission.


[watchmod]: https://github.com/Cyphrme/watchmod
[ETag]:     https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/ETag
[path]:     https://github.com/Cyphrme/Path
