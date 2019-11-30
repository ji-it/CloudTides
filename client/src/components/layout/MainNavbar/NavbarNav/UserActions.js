import React from "react";
import {Link} from "react-router-dom";
import {
    Dropdown,
    DropdownToggle,
    DropdownMenu,
    DropdownItem,
    Collapse,
    NavItem,
    NavLink
} from "shards-react";

import auth from "../../../../utils/auth";

export default class UserActions extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            visible: false
        };

        this.toggleUserActions = this.toggleUserActions.bind(this);
    }

    toggleUserActions() {
        this.setState({
            visible: !this.state.visible
        });
    }

    logout() {
        auth.clearToken();
        auth.clearUserInfo();
        this.props.history.push("/")
    }

    render() {
        return (
            <NavItem tag={Dropdown} toggle={this.toggleUserActions}>
                <DropdownToggle tag={NavLink} className="text-nowrap px-3">
                    <img
                        className="user-avatar rounded-circle mr-2"
                        src={require("./../../../../images/avatars/0.png")}
                        alt="User Avatar"
                    />{" "}
                    {/*<span className="d-none d-md-inline-block">Sierra Brooks</span>*/}
                </DropdownToggle>
                <Collapse tag={DropdownMenu} right small open={this.state.visible}>
                    <DropdownItem tag={Link} to="user-profile">
                        {/*<i className="material-icons">&#xE7FD;</i> */}
                        Profile
                    </DropdownItem>
                    <DropdownItem divider/>
                    <DropdownItem tag={Link} to="#" onClick={this.logout} className="text-danger">
                        {/*<i className="material-icons text-danger">&#xE879;</i> */}
                        Logout
                    </DropdownItem>
                </Collapse>
            </NavItem>
        );
    }
}
