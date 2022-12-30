#!/usr/bin/env bash
# See https://github.com/zamicol/watch
(cd test/watch/src && esbuild test_1.js --format=esm --platform=browser --minify --sourcemap --outfile=test_1?fv=00000000.min.js)


# Calling FileVer here is not needed because FileVer's tests call watch.
# Normally watch must call FileVer, instead of FilveVer calling watch as in this
# example. 
# # Run FileVer from Go test.  In production, this may be done by calling a HTTP
# # endpoint.  
# (cd ../ && go test -run Example_watchVersionAndReplace)









# If using an EndVer format, the following changes will need to be added:   
# Append the map URL to the map file manually since esbuild will always misname
# it. 
#(cd test/watch/src && echo '//# sourceMappingURL=test_1.min.js.map?fv=00000000' >> test_1.min.js)
# Dummy version the min and map file.  (esbuild doesn't allow naming of query
# parameters in names)
#(cd test/watch/src && mv test_1.min.js test_1.min.js?fv=00000000  && mv test_1.min.js.map test_1.min.js.map?fv=00000000 )
