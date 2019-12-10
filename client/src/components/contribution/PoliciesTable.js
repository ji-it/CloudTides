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
import ModalAddPolicy from "../contribution/ModalAddPolicy";
import SuccessNotificationModal from "../common/SuccessNotificationModal";
import ModalAddResource from "../common/ModalAddResource";
import Store from "../../flux/store";
import {Actions} from "../../flux";

class PoliciesTable extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            tableData: Store.getPoliciesTableData()
        };

        this.onChange = this.onChange.bind(this);
    }

    componentWillMount() {
        Store.addChangeListener(this.onChange);
    }

    componentWillUnmount() {
        Store.removeChangeListener(this.onChange);
        clearInterval(this.timer1);
        this.timer1 = null;
    }

    componentDidMount() {
        Actions.getPolicies();
        this.timer1 = Actions.getPolicies(true);
    }

    onChange() {
        this.setState({
            ...this.state,
            tableData: Store.getPoliciesTableData()
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
                                <span className="text text-uppercase">Add Policy</span>
                            </Button>
                            <ModalAddPolicy onExit={this.myCallBack}
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

PoliciesTable.propTypes = {
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

PoliciesTable.defaultProps = {
    title: "Contribution Policies",
    columns: [
        {
            label: "Name",
            name: "name",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    return <b>{value}</b>;
                }
            }
        },
        {
            name: "date_created",
            label: "Date Created",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    const date = new Date(value);
                    return date.toUTCString();
                }
            }
        },
        {
            name: "project_name",
            label: "Project",
        },
        {
            label: "Deploy Type",
            name: "deploy_type",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    return <b>{value}</b>;
                }
            }
        },
        {
            label: "Idle %",
            name: "idle_policy",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    const json = JSON.parse(value);
                    console.log(json)
                    const nf = new Intl.NumberFormat('en-US', {
                        style: 'percent',
                    });
                    return (
                        "CPU = " + nf.format(json.cpu) + ", RAM = " + nf.format(json.ram)
                    );
                }
            }
        },
        {
            label: "Threshold %",
            name: "threshold_policy",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    const json = JSON.parse(value);
                    const nf = new Intl.NumberFormat('en-US', {
                        style: 'percent',
                    });
                    return (
                        "CPU = " + nf.format(json.cpu) + ", RAM = " + nf.format(json.ram)
                    );
                }
            }
        },
        {
            name: "is_destroy",
            label: "Destroy/Stop",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    return (value) ? "Destroy" : "Stop"
                }
            }
        },
        {
            name: "hosts_assigned",
            label: "Hosts Assigned",
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

export default PoliciesTable;
