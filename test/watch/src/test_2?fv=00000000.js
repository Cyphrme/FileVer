import * as test1 from './test_1?fv=00000000.min.js';
import * as test3 from './subdir/test_3?fv=00000000.js';
import * as test4 from './subdir/test_4?fv=00000000.js';
// Comments referring to './test_1?fv=00000000.min.js' should be updated as
// well, but comments referring to `test_1.js` or './test_1?fv=00000000.js' will
// be left untouched.