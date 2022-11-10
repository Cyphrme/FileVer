#!/usr/bin/env bash
#
# See https://github.com/zamicol/watch
#
# For prod, use `min.js`. For development, `join.js` may be useful. 
(cd test/watch/src && esbuild test_1.js --format=esm --platform=browser --minify --sourcemap=external --outfile=test_1.min.js)

# Append the map URL to the map file manually since esbuild will always misname
# it. 
(cd test/watch/src && echo '//# sourceMappingURL=test_1.min.js.map?fv=00000000' >> test_1.min.js)

# Dummy version the min and map file.  (esbuild doesn't allow naming of query
# parameters in names)
(cd test/watch/src && mv test_1.min.js test_1.min.js?fv=00000000  && mv test_1.min.js.map test_1.min.js.map?fv=00000000 )

# Calling FileVer here is not needed because FileVer's tests call watch.
# Normally watch must call FileVer.  
# # Run FileVer from Go test.  In production, this may be done by calling a HTTP
# # endpoint.  
# (cd ../ && go test -run Example_watchVersionAndReplace)
