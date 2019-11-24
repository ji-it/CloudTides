import React from "react";
import PropTypes from "prop-types";
import {Navbar, NavbarBrand} from "shards-react";
import Actions from '../../../flux/actions';


class SidebarMainNavbar extends React.Component {
    constructor(props) {
        super(props);

        this.handleToggleSidebar = this.handleToggleSidebar.bind(this);
    }

    handleToggleSidebar() {
        Actions.toggleMenu();
    }

    render() {
        const {hideLogoText} = this.props;
        return (
            <div className="main-navbar">
                <Navbar
                    className="align-items-stretch bg-white flex-md-nowrap border-bottom p-0"
                    type="light"
                >
                    <NavbarBrand
                        className="w-100 mr-0"
                        href="#"
                        style={{lineHeight: "25px"}}
                    >
                        <div className="d-table ml-30">
                            <img
                                id="main-logo"
                                className="d-inline-block align-top mr-1"
                                style={{maxWidth: "65px"}}
                                src={require("../../../images/logo/tides logo_dark.png")}
                                alt="Clouds Tides Logo"
                            />
                            {/*            {!hideLogoText && (*/}
                            {/*                <span className="d-none d-md-inline ml-1">*/}
                            {/*  Shards Dashboard*/}
                            {/*</span>*/}
                            {/*            )}*/}
                        </div>
                    </NavbarBrand>
                    {/* eslint-disable-next-line */}
                    <a
                        className="toggle-sidebar d-sm-inline d-md-none d-lg-none"
                        onClick={this.handleToggleSidebar}
                    >
                        <i className="material-icons">&#xE5C4;</i>
                    </a>
                </Navbar>
            </div>
        );
    }
}

SidebarMainNavbar.propTypes = {
    /**
     * Whether to hide the logo text, or not.
     */
    hideLogoText: PropTypes.bool
};

SidebarMainNavbar.defaultProps = {
    hideLogoText: false
};

export default SidebarMainNavbar;
