import React from "react";
import PropTypes from "prop-types";
import {
    Row,
    Col,
    Card,
    CardHeader,
    CardBody,
    Progress,
} from "shards-react";

class WorkloadsCard extends React.Component {

    render() {
        const {title} = this.props;
        return (
            <Card small className="h-100">
                <CardHeader className="border-bottom">
                    <h6 className="m-0">{title}</h6>
                </CardHeader>
                <CardBody className="d-flex py-0">
                    <div className="w-100 mx-10 text-left">
                        <div className="mt-2" style={{fontSize: "1.2em", color: "#1B2376"}}>
                            <span style={{fontSize: "1.8em"}}>400 jobs</span> contributed
                        </div>
                        <div className="mb-3 font-weight-normal">
                            <div><b>30</b> running | <b>4</b> suspended</div>
                            <div><b>70</b> resources used</div>
                        </div>
                        <div className="mb-3">
                            <Progress className="workloadProgress" barClassName="workloadProgressBar" value="20">20</Progress>
                        </div>
                    </div>
                </CardBody>
            </Card>
        );
    }
}

WorkloadsCard.propTypes = {
    /**
     * The component's title.
     */
    title: PropTypes.string,
};

WorkloadsCard.defaultProps = {
    title: "Workloads",
};

export default WorkloadsCard;
