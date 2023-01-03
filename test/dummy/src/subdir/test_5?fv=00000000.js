import * as test1 from '../test_1?fv=00000000.js'; //"Relative in parent dir"
import * as test2 from '../test_2?fv=00000000.js'; //"Relative in parent dir"
import * as test3 from '../subdir/test_3?fv=00000000.js'; // "Relative in current dir from root".  
import * as test3 from './test_3?fv=00000000.js'; // "Relative in current subdir" // TODO does not work.  