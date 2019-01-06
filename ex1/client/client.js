'use strict';

var fc = require('fabric-client');
var config = require('./config.js')

var channel = fc.newChannel(config.CHANNEL_NAME);

function invoke(chaincodeID, args) {
    
}