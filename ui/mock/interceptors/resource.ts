import { commerce, random, address, internet } from 'faker';


export const postHandler = (req, res, next) => {
  req.body = {
    name: address.city(),
    host: internet.url(),
    status: commerce.productName(),
    sites: random.number(10),
    orgs: random.number(10),
    vms: random.number(10),
    services: random.number(10),
  }
  next();
}