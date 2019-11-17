import React from "react";
import PropTypes from "prop-types";
import {NavLink as RouteNavLink} from "react-router-dom";
import {NavItem, NavLink} from "shards-react";

const ChildSidebarNavItem = ({item}) => (
    <NavItem>
        <NavLink className="child" tag={RouteNavLink} to={item.to}>
            {item.name && <span>{item.name}</span>}
        </NavLink>
    </NavItem>
);

ChildSidebarNavItem.propTypes = {
    /**
     * The item object.
     */
    item: PropTypes.object
};

export default ChildSidebarNavItem;
