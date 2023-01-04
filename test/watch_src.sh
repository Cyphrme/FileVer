#!/usr/bin/env bash
# See https://github.com/zamicol/watch

# For this example, the map is not versioned since it would conflict with the
# min file name, e.g. regex `test_1\?fv=[0-9A-Za-z_-]*.min.js` matches
# `test_1?fv=00000000.min.js` and `test_1?fv=00000000.min.js.map`.  There are a
# few ways of handling this, including naming the map file
# `test_1?fv=00000000.map.min.js` instead.  Also, browsers are not suppose to
# cache maps files, so versioning doesn't appear useful.  
# For this example, the map file is not versioned.  
(cd test/watch/src && esbuild test_1.js --format=esm --platform=browser --minify --sourcemap --outfile=test_1.min.js)
# Do not version the map file, but add dummy version to the min.  (Esbuild
# doesn't allow naming maps, so the min file has to be renamed after
# generation).  
(cd test/watch/src && mv test_1.min.js.map ../dist/test_1.min.js.map && mv test_1.min.js test_1?fv=00000000.min.js)


# Calling FileVer here is not needed because FileVer's tests call watch.
# Normally watch must call FileVer, instead of FilveVer calling watch as in this
# example. 
# # Run FileVer from Go test.  In production, this may be done by calling a HTTP
# # endpoint.  
# (cd ../ && go test -run Example_watchVersionAndReplace)








# Alternativly, the mapp information can be manually appeneded to the end of the file:
# If using an EndVer format, the following changes will need to be added:   
# Append the map URL to the map file manually since esbuild will always misname
# it. 
#(cd test/watch/src && echo '//# sourceMappingURL=test_1.min.js.map?fv=00000000' >> test_1.min.js)

