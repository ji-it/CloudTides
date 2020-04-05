import { helpers } from 'faker';

const users = [
  'admin@vmware.com',
  'andrewz@vmware.com',
]

export const postLogin = (req, res, next) => {
  const { username, password } = req.body;
  if (password === "test" && users.includes(username)) {
    const item = helpers.userCard()
    item.username = username
    req.body = {
      ...item,
      admin: username === users[0]
    };
  } else {
    req.status = 401;
  }
  next();
}