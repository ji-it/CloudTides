import React from 'react';
import {makeStyles} from '@material-ui/core/styles';
import TreeView from '@material-ui/lab/TreeView';
import FiberManualRecordIcon from '@material-ui/icons/FiberManualRecord';
import ArrowDropDownIcon from '@material-ui/icons/ArrowDropDown';
import ArrowRightIcon from '@material-ui/icons/ArrowRight';
import TreeItem from '@material-ui/lab/TreeItem';
import Actions from '../../../flux/actions'

const useStyles = makeStyles({
    root: {
        height: 216,
        flexGrow: 1,
        maxWidth: 400,
    },
});

export default function TreeChildbarNavItem({data}) {
    const classes = useStyles();
    const [expanded, setExpanded] = React.useState([]);

    const handleChange = (event, nodes) => {
        setExpanded(nodes);
    };

    const handleClick = (event) => {
        const idx = event.target.name;
        Actions.switchResource(idx)
    };

    return (
        <TreeView
            className={classes.root}
            defaultCollapseIcon={<ArrowDropDownIcon style={{color: "#4875B3", fontSize: "20px"}}/>}
            defaultExpandIcon={<ArrowRightIcon style={{color: "#4875B3", fontSize: "20px"}}/>}
            defaultEndIcon={<FiberManualRecordIcon style={{color: "#9BC225", fontSize: "10px"}}/>}
            expanded={expanded}
            onNodeToggle={handleChange}
        >
            {data && data.map((item, idx) => {
                return (<a name={idx} key={idx} onClick={handleClick} className={"ml-2"}
                           style={{cursor: "Pointer"}}
                >
                    {item.datacenter}
                </a>)
            })}
        </TreeView>
    );
}