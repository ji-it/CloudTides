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

class Resources extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            tableData: Store.getResourceTableData()
        };

        this.onChange = this.onChange.bind(this);
    }

    componentWillMount() {
        Store.addChangeListener(this.onChange);
    }

    componentWillUnmount() {
        Store.removeChangeListener(this.onChange);
    }

    componentDidMount() {
        Actions.getResources()
    }

    onChange() {
        this.setState({
            ...this.state,
            tableData: Store.getResourceTableData()
        });
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
        const {title, columns, options} = this.props;
        const {tableData: data} = this.state;
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
                            <ModalAddResource onExit={this.myCallBack}
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
    tableStyle: PropTypes.object,
    options: PropTypes.object,
};

Resources.defaultProps = {
    title: "Resources",
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
            name: "ip_address",
            label: "IP Address"
        },
        {
            name: "cpu",
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
                            {
                                (value == 0.4) ?
                                    <ArrowDropUpIcon className="text-success"/>
                                    :
                                    <ArrowDropDownIcon className="text-danger"/>
                            }
                        </div>
                    );
                }
            }
        },
        {
            name: "memory",
            label: "Memory",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    return (
                        <div>
                            {value}
                            {
                                (value == 0.4) ?
                                    <ArrowDropUpIcon className="text-success"/>
                                    :
                                    <ArrowDropDownIcon className="text-danger"/>
                            }
                        </div>
                    );
                }
            }
        },
        {
            name: "disk",
            label: "Disk",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    return (
                        <div>
                            {value}
                            {
                                (value == 0.4) ?
                                    <ArrowDropUpIcon className="text-success"/>
                                    :
                                    <ArrowDropDownIcon className="text-danger"/>
                            }
                        </div>
                    );
                }
            }
        },
        {
            name: "jobs_done",
            label: "Jobs Done"
        },
        {
            name: "policy",
            label: "Policy",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    return <b>{value}</b>;
                }
            }
        },
        {
            name: "active",
            label: "Active",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    return (
                        <FormControlLabel
                            label={value ? "Yes" : "No"}
                            value={value ? "Yes" : "No"}
                            control={
                                <Switch color="primary" checked={value} value={value ? "Yes" : "No"}/>
                            }
                            onChange={event => {
                                updateValue(event.target.value === "Yes" ? false : true);
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
