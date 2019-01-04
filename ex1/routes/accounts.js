'use strict';
var express = require('express');
var router = express.Router();

var fc = require('fabric-client');

var mychannel = fc.newChannel("mychannel");

/* GET accounts listing. */
router.get('/', function(req, res, next) {
  res.send('respond with a resource');
});

router.post('/api/create-account', async (req, res, next) => {
  req.
  let invokeReq = {
    chaincodeId: 'accountmgmt',
    fcn: 'CreateAccount',
    args: {
      
    }
  };
  await mychannel.sendTransactionProposal();
});
module.exports = router;
