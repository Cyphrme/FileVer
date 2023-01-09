import * as test1 from '../test_1~fv=00000000.js'; // "Relative in parent dir"
import * as test2 from '../test_2~fv=00000000.js'; // "Relative in parent dir"
import * as test3 from '../subdir/test_3~fv=00000000.js'; // "Relative in current dir from root".  
// "Relative in current subdirectory" **Does not work**.  References must be always relative to root.  See README.  
import * as test3 from './test_3~fv=00000000.js'; 