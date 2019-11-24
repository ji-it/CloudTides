import React from "react";
import PropTypes from "prop-types";
import classNames from "classnames";
import {Col} from "shards-react";

import SidebarMainNavbar from "./SidebarMainNavbar";
import LeftSidebarNavItems from "./LeftSidebarNavItems";
import RightSidebarNavItems from "./RightSidebarNavItems";
import Store from "../../../flux/store";

class MainSidebar extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            menuVisible: false,
            sidebarNavItems: Store.getSidebarItems()
        };

        this.onChange = this.onChange.bind(this);
    }

    componentWillMount() {
        Store.addChangeListener(this.onChange);
    }

    componentWillUnmount() {
        Store.removeChangeListener(this.onChange);
    }

    onChange() {
        this.setState({
            ...this.state,
            menuVisible: Store.getMenuState(),
            sidebarNavItems: Store.getSidebarItems()
        });
    }

    render() {
        const classes = classNames(
            "main-sidebar",
            "px-0",
            // "col-12",
            this.state.menuVisible && "open"
        );
        const {sidebarNavItems: items} = this.state;
        let hasChildren = items.filter(
            item => item.to === this.props.location.pathname
        ).shift();
        hasChildren = hasChildren && hasChildren.show;
        return (
            <Col
                tag="aside"
                className={classes}
                // lg={{size: 2}}
                // md={{size: 3}}
            >
                <SidebarMainNavbar hideLogoText={this.props.hideLogoText}/>
                {/*<SidebarSearch />*/}
                <Col
                    className="p-0 left-bar"
                >
                    <LeftSidebarNavItems/>
                </Col>
                {hasChildren ? (
                    <Col
                        className="p-0 right-bar"
                    >
                        <RightSidebarNavItems {...this.props}/>
                    </Col>
                ) : (
                    <div></div>
                )}
            </Col>
        );
    }
}

MainSidebar.propTypes = {
    /**
     * Whether to hide the logo text, or not.
     */
    hideLogoText: PropTypes.bool
};

MainSidebar.defaultProps = {
    hideLogoText: false
};

export default MainSidebar;
