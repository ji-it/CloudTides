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
import ModalAddResource from "../common/ModalAddResource";
import SuccessNotificationModal from "../common/SuccessNotificationModal";
import Store from "../../flux/store";
import {Actions} from "../../flux";
import {devURL} from "../../utils/urls";
import request from "../../utils/request";
import auth from "../../utils/auth";

class Resources extends React.Component {
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
        const {title, columns, options, data, policiesData} = this.props;
        return (
            <Card small className="blog-comments">
                <CardHeader className="m-2 mb-0">
                    <div>
                        <div style={{display: "inline-block"}}>
                            <h6 className="m-0 font-weight-bold">{title}</h6>
                        </div>
                        <div style={{display: "inline-block", float: "right"}}>
                            <Button
                                className="shadow-sm"
                                onClick={() => this.toggleModal("addModal")}
                            >
                                <span className="text text-uppercase">Add Resource</span>
                            </Button>
                            <ModalAddResource policiesData={policiesData} onExit={this.myCallBack}
                                              toggleState={this.state.addModal}/>
                            <SuccessNotificationModal onRef={ref => (this.successmodal = ref)}/>
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

Resources.propTypes = {
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

Resources.defaultProps = {
    title: "Resources",
    data: [],
    columns: [
        {
            name: "id",
            label: "ID",
            options: {
                display: false
            }
        },
        {
            name: "host_name",
            label: "Name",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    return <b>{value}</b>;
                }
            }
        },
        {
            name: "status",
            label: "Status",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    const style = "tides-" + value.toLowerCase();
                    return (
                        <span className={style + " border small status-border"}
                              data-task-status={value.toLowerCase()}
                        >
                            {value}
                        </span>
                    );
                }
            }
        },
        {
            name: "host_name",
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
            name: "job_completed",
            label: "Jobs Done"
        },
        {
            name: "policy_name",
            label: "Policy",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    return <b>{value}</b>;
                }
            }
        },
        {
            name: "is_active",
            label: "Active",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    return (
                        <FormControlLabel
                            label={value ? "Yes" : "No"}
                            value={value ? "Yes" : "No"}
                            name={tableMeta.rowData[0]}
                            control={
                                <Switch color="primary" checked={value} value={value ? "Yes" : "No"}/>
                            }
                            onChange={event => {
                                const id = event.target.name;
                                updateValue(event.target.value === "Yes" ? false : true);
                                const endpoint = '/api/resource/toggle_active/';
                                const formData = {id: id};
                                const requestURL = devURL + endpoint;
                                request(requestURL, {method: 'POST', body: formData})
                                    .then((response) => {
                                        //Load dashboard data:- resource list, total contribution (cost and power), total resource usage (usage use + hosts number, idle, vms)
                                        // this.redirectUser();
                                    }).catch((err) => {
                                    console.log(err);
                                });
                            }}
                        />
                    );

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

export default Resources;
