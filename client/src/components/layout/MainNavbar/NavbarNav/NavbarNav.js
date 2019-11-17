import React from "react";
import { Nav } from "shards-react";

import Notifications from "./Notifications";
import UserActions from "./UserActions";

export default () => (
  <Nav navbar className="flex-row ml-auto">
    {/*<Notifications />*/}
    <UserActions />
  </Nav>
);
