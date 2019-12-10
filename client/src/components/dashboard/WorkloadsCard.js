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
        const {title, data} = this.props;
        const {workload} = data;
        let contributed = 0, running = 0, destroyed = 0, hosts_used = 0, percentage;
        if (workload) {
            contributed = workload.contributed;
            running = workload.running;
            destroyed = workload.destroyed;
            hosts_used = workload.hosts_used;
            percentage = running / contributed * 100;
        }
        const nf = new Intl.NumberFormat('en-US', {
            maximumFractionDigits: 2,
        });
        return (
            <Card small className="h-100">
                <CardHeader className="border-bottom">
                    <h6 className="m-0">{title}</h6>
                </CardHeader>
                <CardBody className="d-flex py-0">
                    <div className="w-100 mx-10 text-left">
                        <div className="mt-2" style={{fontSize: "1.2em", color: "#1B2376"}}>
                            <span style={{fontSize: "1.8em"}}>{contributed}</span> contributed
                        </div>
                        <div className="mb-3 font-weight-normal">
                            <div><b>{running}</b> running | <b>{destroyed}</b> destroyed</div>
                            <div><b>{hosts_used}</b> resources used</div>
                        </div>
                        <div className="mb-3">
                            <Progress className="workloadProgress" barClassName="workloadProgressBar"
                                      value={String(percentage)}>{nf.format(percentage)}</Progress>
                        </div>
                    </div>
                </CardBody>
            </Card>
        );
    }
}

WorkloadsCard
    .propTypes = {
    /**
     * The component's title.
     */
    title: PropTypes.string,
};

WorkloadsCard
    .defaultProps = {
    title: "Workloads",
};

export default WorkloadsCard;
