'use strict';

const Client = require('fabric-client');
const fs = require('fs');
const path = require('path');

Client.addConfigFile(path.join(__dirname, './config.json'));
var NETWORK = Client.getConfigSetting('network');
var client = new Client();
var channel = client.newChannel(NETWORK['channel-name']);

channel.addOrderer(client.newOrderer(
    NETWORK.orderer.url,
    {
        'pem': Buffer.from(fs.readFileSync(path.join(__dirname, NETWORK.orderer.tls_cacerts))).toString(),
        'ssl-target-name-override': NETWORK.orderer['server-hostname']
    }
));

for (let peer in NETWORK) {
    if (NETWORK.hasOwnProperty(peer) && typeof NETWORK[peer].peer1 !== 'undefined') {
        channel.addPeer(client.newPeer(
            NETWORK[peer].peer1.requests,
            {
                'pem': Buffer.from(fs.readFileSync(path.join(__dirname, NETWORK[peer].peer1.tls_cacerts))).toString(),
                'ssl-target-name-override': NETWORK[peer].peer1['server-hostname']
            }
        ));
    }
}

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
        let errStr = [];
        proposalResps.forEach(resp => {
            let pass = false;
            if (!resp.response) {
                if (resp instanceof Error) {
                    errStr.push(resp.message);
                } else {
                    errStr.push('resp.response null');
                }
            } else if (resp.response.status == 200) {
                pass = channel.verifyProposalResponse(resp);
                if (!pass) {
                    console.log('verifyProposalResponse failed: %o', resp);
                    errStr.push('verify endorser signature error');
                }
            } else {
                console.log('proposal response status: %d', resp.response.status);
                errStr.push('proposal response status: ' + resp.response.status);
            }
            allPass = allPass & pass;
        });
        if (!allPass) {
            throw new Error('verify endorser response: %', errStr);
        }
        allPass = channel.compareProposalResponseResults(proposalResps);
        if (!allPass) {
            throw new Error('endorsing results are different');
        }

        let eventPromises = channelEventHubs.map((eventHub) => {
            return new Promise((resolve) => {
                eventHub.registerTxEvent(txId.getTransactionID(), (txIdStr, code) => {
                    eventHub.unregisterTxEvent(txIdStr);
                    if (code != 'VALID') {
                        throw new Error('peer ledger commit error: ' + code);
                    }
                });
                eventHub.connect();
                resolve();
            });
        });
        return Promise.all(eventPromises).then( () => {
            const txReq = {
                proposalResponses: proposalResps,
                proposal: proposal,
            };
            return channel.sendTransaction(txReq);
        });
    }).then((result) => {
        if (result && result.status === '200') {
        
        } else {
            throw new Error(util.format('sendTransaction error: %o', result));
        }
    });
}

function query(chaincodeId, fcn, args) {
    const request = {
        chaincodeId: chaincodeId,
        fcn: fcn,
        args: args
    };
    return channel.queryByChaincode(request);
}

module.exports.invoke = invoke;
module.exports.query = query; 