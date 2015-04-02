/**
 * Created by languid on 4/2/15.
 */
define(['jquery'], function($){
    if(_inlineCodes && _inlineCodes.length){
        $.map(_inlineCodes, function(fn){
            typeof fn === 'function' && fn()
        })
    }
});