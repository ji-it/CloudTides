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

class TemplatesTable extends React.Component {
    state = {};
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
        const {title, columns, data, options} = this.props;
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
                                <span className="text text-uppercase">Upload File</span>
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

TemplatesTable.propTypes = {
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

TemplatesTable.defaultProps = {
    title: "Templates Manager",
    columns: [
        {
            name: "Name",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    return <b>{value}</b>;
                }
            }
        },
        "Date Added",
        {
            name: "Guest OS",
            options: {
                filter: true,
                customBodyRender: (value, tableMeta, updateValue) => {
                    return <b>{value}</b>;
                }
            }
        },
        "Compatibility", "Provisioned Space", "Memory Size",

        // {name: "", options: { filter: false,sort: false,empty: true,customBodyRender: (value, tableMeta, updateValue) => {return ();}}},
    ],
    data: [
        ["kube-1107", "26/09/2019", "Ubuntu Linux (64-bit)", "ESXI 6.5 and later", "34.2GB", "1GB"],
        ["kube-1108", "26/09/2019", "Ubuntu Linux (64-bit)", "ESXI 6.5 and later", "34.2GB", "1GB"],
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

export default TemplatesTable;
