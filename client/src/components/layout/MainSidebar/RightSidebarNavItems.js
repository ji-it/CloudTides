import React from "react";
import {Nav} from "shards-react";

import ChildSidebarNavItem from "./ChildSidebarNavItem";
import {Store} from "../../../flux";

class RightSidebarNavItems extends React.Component {
    constructor(props) {
        super(props)

        this.state = {
            navItems: Store.getSidebarTitles()
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
            navItems: Store.getSidebarTitles()
        });
    }

    render() {
        const {navItems: items} = this.state;
        const selectedItem = items.filter(item => item.to === this.props.location.pathname).shift();
        console.log(selectedItem);
        return (
            <div className="child-nav-wrapper">
                <h6 className="ml-10 mt-10 font-weight-bold text-black">{selectedItem.name}</h6>
                <hr className="mx-10"/>
                <Nav className="nav--no-borders flex-column">
                    {selectedItem.children && selectedItem.children.map((item, idx) => (
                        <ChildSidebarNavItem key={idx} item={item}/>
                    ))}
                </Nav>
            </div>
        )
    }
}

export default RightSidebarNavItems;
