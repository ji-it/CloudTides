import React from 'react';
import {makeStyles} from '@material-ui/core/styles';
import TreeView from '@material-ui/lab/TreeView';
import FiberManualRecordIcon from '@material-ui/icons/FiberManualRecord';
import ArrowDropDownIcon from '@material-ui/icons/ArrowDropDown';
import ArrowRightIcon from '@material-ui/icons/ArrowRight';
import TreeItem from '@material-ui/lab/TreeItem';

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
                const dt = new Date().getTime();
                return (<TreeItem key={idx} nodeId={String(idx)}
                                  label={item.name}>
                    {item.children && item.children.map((item, idx) => (
                        <TreeItem key={idx} nodeId={item.name + String(idx + dt + Math.random())}
                                  label={item.name}/>
                    ))}
                </TreeItem>)
            })}
        </TreeView>
    );
}