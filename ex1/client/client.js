'use strict';

const Client = require('fabric-client');
const config = require('./config.js')

var client = new Client();
var channel = client.newChannel(config.CHANNEL_NAME);
channel.addOrderer();
channel.addPeer();
channel.initialize();
var channelEventHubs = channel.getPeers().map((channelPeer) => {
    return channelPeer.getChannelEventHub();
});

function invoke(chaincodeId, fcn, args) {
    const txId = client.newTransactionID();
    const req = {
        chaincodeId: chaincodeId,
        fcn: fcn,
        args: args,
        txId: txId,
    };
    return channel.sendTransactionProposal(req).then((results) => {
        const proposalResps = results[0];
        const proposal = results[1];
        let allPass = true;
        
        proposalResps.forEach(resp => {
            let pass = false;
            if (resp.response && resp.response.status == 200) {
                pass = channel.verifyProposalResponse(resp);
                if (pass === false) {

                }
            } else {

            }
            allPass = allPass & pass;
        });
        if (allPass === false) {

        }
        allPass = channel.compareProposalResponseResults(proposalResps);
        if (allPass === false) {

        }

        let eventPromises = channelEventHubs.map((eventHub) => {
            return new Promise((resolve) => {
                eventHub.registerTxEvent(txId.getTransactionID(), (txIdStr, code) => {
                    eventHub.unregisterTxEvent(txIdStr);
                    if (code != 'VALID') {
                        
                    }
                });
                eventHub.connect();
            });
        });
        return Promise.all(eventPromises).then( => {
            return {
                proposalResponses: proposalResps,
                proposal: proposal,
            };
        },
        (err) => {

        });
    }).then((txReq)=> {
            channel.sendTransaction(txReq);
    }).then((result) => {
        if (result && result.status === "200") {
        }
    },
    (err) => {

    });
}