'use strict';
var express = require('express');
var router = express.Router();

var client = require('../client/client.js');

/* GET accounts listing. */
router.get('/', async function(req, res, next) {
  try {
    result = await client.query('accountmgmt', 'QueryIDByIDOrName', ['xccc']);
    res.send(result);
  } catch(e) {
    res.send('query error:' + e);
  }
});

router.post('/api/create-account', async (req, res, next) => {
  try {
    await client.invoke('accountmgmt', 'CreateAccount', ['a', 'b', 1000]);
    res.send('invoke success');
  } catch(e) {
    res.send('invoke error: ' + e);
  }
});

module.exports = router;
