import React from "react";
import PropTypes from "prop-types";
import {
    Card,
    CardHeader,
    CardBody,
    Button,
} from "shards-react";
import MUIDataTable from "mui-datatables"
import {createMuiTheme, MuiThemeProvider} from '@material-ui/core/styles';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Switch from '@material-ui/core/Switch';
import ArrowDropUpIcon from '@material-ui/icons/ArrowDropUp';
import ArrowDropDownIcon from '@material-ui/icons/ArrowDropDown';
import SuccessNotificationModal from "../common/SuccessNotificationModal";
import Store from "../../flux/store";
import {Actions} from "../../flux";

class VMTable extends React.Component {
    constructor(props) {
        super(props);

        this.state = {};
    }

    toggleModal = state => {
        this.setState({
            [state]: !this.state[state]
        });
    };

    myCallBack = (data, showNotif) => {
        this.toggleModal(data);
        if (showNotif)
            this.successmodal.toggleModal("notificationModal");
    };

    getMuiTheme = () => createMuiTheme({
        overrides: this.props.tableStyle
    });

    render() {
        const {title, columns, options, data} = this.props;
        return (
            <Card small className="blog-comments">
                <CardHeader className="m-2 mb-0">
                    <div>
                        <div style={{display: "inline-block"}}>
                            <h6 className="m-0 font-weight-bold">{title}</h6>
                        </div>
                    </div>
                </CardHeader>
                <CardBody className="p-0">
                    <MuiThemeProvider theme={this.getMuiTheme()}>
                        <MUIDataTable
                            data={data}
                            columns={columns}
                            options={options}
                        />
                    </MuiThemeProvider>
                </CardBody>
            </Card>
        );
    }
}

VMTable.propTypes = {
    /**1
     * The component's title.
     */
    title: PropTypes.string,
    /**
     * The table dataset.
     */
    columns: PropTypes.array,
    data: PropTypes.array,
    tableStyle: PropTypes.object,
    options: PropTypes.object,
};

VMTable.defaultProps = {
    title: "Virtual Machines",
    data: [],
    columns: [
        {
            name: "name",
            label: "Name",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    return <b>{value}</b>;
                }
            }
        },
        {
            name: "powered_on",
            label: "Status",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    const style = value ? "primary" : "warning";
                    return (
                        <span className={style}
                        >
                            {value ? "On" : "Off"}
                        </span>
                    );
                }
            }
        },
        {
            name: "ip_address",
            label: "IP Address"
        },
        {
            name: "cpu_percent",
            label: "CPU",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    const nf = new Intl.NumberFormat('en-US', {
                        style: 'percent',
                    });
                    return (
                        <div>
                            {nf.format(value)}
                            {/*{*/}
                            {/*    (value == 0.4) ?*/}
                            {/*        <ArrowDropUpIcon className="text-success"/>*/}
                            {/*        :*/}
                            {/*        <ArrowDropDownIcon className="text-danger"/>*/}
                            {/*}*/}
                        </div>
                    );
                }
            }
        },
        {
            name: "ram_percent",
            label: "Memory",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    const nf = new Intl.NumberFormat('en-US', {
                        style: 'percent',
                    });
                    return (
                        <div>
                            {nf.format(value)}
                            {/*{*/}
                            {/*    (value == 0.4) ?*/}
                            {/*        <ArrowDropUpIcon className="text-success"/>*/}
                            {/*        :*/}
                            {/*        <ArrowDropDownIcon className="text-danger"/>*/}
                            {/*}*/}
                        </div>
                    );
                }
            }
        },
        // {
        //     name: "total_disk",
        //     label: "Disk",
        //     options: {
        //         filter: true,
        //         customBodyRender: (value, tableMeta, updateValue) => {
        //             return (
        //                 <div>
        //                     {value}
        //                     {
        //                         (value == 0.4) ?
        //                             <ArrowDropUpIcon className="text-success"/>
        //                             :
        //                             <ArrowDropDownIcon className="text-danger"/>
        //                     }
        //                 </div>
        //             );
        //         }
        //     }
        // },
        {
            name: "guest_os",
            label: "Guest OS"
        },
        {
            name: "date_created",
            label: "Added",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    const date = new Date(value);
                    const res = (value) ? date.toUTCString() : "";
                    return res;
                }
            }
        },
        {
            name: "date_destroyed",
            label: "Destroyed",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    const date = new Date(value);
                    const res = (value) ? date.toUTCString() : "";
                    return res;
                }
            }
        },
        {
            name: "boinc_time",
            label: "BOINC Start",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    const date = new Date(value);
                    const res = (value) ? date.toUTCString() : "";
                    return res;
                }
            }
        },
        // {name: "", options: { filter: false,sort: false,empty: true,customBodyRender: (value, tableMeta, updateValue) => {return ();}}},
    ],
    tableStyle: {
        MUIDataTableSelectCell: {
            headerCell: {
                backgroundColor: "#E9EDF6",
            }
        },
        MUIDataTableHeadCell: {
            fixedHeader: {
                backgroundColor: "#E9EDF6",
            },
            data: {
                margin: "auto"
            }
        },
        MUIDataTableBodyCell: {
            root: {
                textAlign: "center",
            }
        }
    },
    options: {
        filterType: 'checkbox',
        customToolbarSelect: () => {
        },
        elevation: 0,
        filter: false,
        responsive: "scrollMaxHeight",
    },
};

export default VMTable;
