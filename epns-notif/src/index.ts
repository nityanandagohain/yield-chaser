import epnsHelper, { InfuraSettings, NetWorkSettings, EPNSSettings } from '@epnsproject/backend-sdk-staging';
import config from "./config/config";
import express, { json } from "express";

let projectId = process.env.PROJECT_ID as string
let projectSecret = process.env.PROJECT_SECRET as string
let privateKey = process.env.PRIVATE_KEY as string

console.log(projectId, projectSecret, privateKey)
// InfuraSettings contains setttings details on infura
const infuraSettings: InfuraSettings = {
  projectID: projectId,
  projectSecret: projectSecret 
}

// Network settings contains details on alchemy, infura and etherscan
const settings: NetWorkSettings = {
  // alchemy: config.alchemyAPI,
  infura: infuraSettings,
  etherscan: config.etherscanAPI
}

// EPNSSettings settings contains details on EPNS network, contract address and contract ABI
const epnsSettings: EPNSSettings = {
  network: config.web3RopstenNetwork,
  contractAddress: config.deployedContract,
  contractABI: config.deployedContractABI
}

let channelPrivateKey = privateKey
const sdk = new epnsHelper(config.web3RopstenNetwork, channelPrivateKey, settings, epnsSettings)

let user = "0xfDdA054f4C5A9bCFC8512f5Cf220E6E77430C556"
let title = "yield-chaser"
let message = ""
const sendNotification = async (payloadTitle :string, payloadMessage :string, notificationType :number) => {
  await sdk.sendNotification(user, title, message, payloadTitle, payloadMessage, notificationType, false)
};




const app = express();
const port = 8080; // default port to listen
app.use(express.json());

// define a route handler for the default home page
app.get("/", (req, res) => {
  res.send("Yield Chaser, notification service using EPNS");
});


app.post("/notification", async (req, res) => {
  let jsonBody = req.body;
  if (("title" in jsonBody) == false){
    res.send("wrong request body, please send titie");
    return
  }
  if (("message" in jsonBody) == false){
    res.send("wrong request body, please send message");
    return
  }
  if (("notificationType" in jsonBody) == false){
    res.send("wrong request body, please send notificationType, 1, 2, 3");
    return
  }

  console.log("sending notification for :", jsonBody)
  await sendNotification(jsonBody["title"], jsonBody["message"], jsonBody["notificationType"])
  res.send("notification sent to the cannels.");
});

// start the Express server
app.listen(port, () => {
  // tslint:disable-next-line:no-console
  console.log(`server started at http://localhost:${port}`);
});

