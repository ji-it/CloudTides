import React from "react";
import {Nav} from "shards-react";

import UserActions from "./UserActions";

export default (props) => (
    <Nav navbar className="flex-row ml-auto">
        <UserActions {...props}/>
    </Nav>
);
