import React from "react";
import PropTypes from "prop-types";
import {
    Card,
    CardHeader,
    Button,
    ListGroup,
    ListGroupItem,
    Progress
} from "shards-react";

class UserDetails extends React.Component {


    render() {
        const {data} = this.props;
        data.avatar = require("./../../images/avatars/0.png");

        return (
            <Card small className="mb-4 pt-3">
                <CardHeader className="border-bottom text-center">
                    <div className="mb-3 mx-auto">
                        <img
                            className="rounded-circle"
                            src={data.avatar}
                            alt={data.avatar}
                            width="110"
                        />
                    </div>
                    <h4 className="mb-0">{data.first_name + ' ' + data.last_name}</h4>
                    <span className="text-muted d-block mb-2">{data.position}</span>
                    {/*<Button pill outline size="sm" className="mb-2">*/}
                    {/*    <i className="material-icons mr-1">edit</i> Save*/}
                    {/*</Button>*/}
                </CardHeader>
                {/*<ListGroup flush>*/}
                {/*    <ListGroupItem className="px-4">*/}
                {/*        <div className="progress-wrapper">*/}
                {/*            <strong className="text-muted d-block mb-2">*/}
                {/*                {userDetails.performanceReportTitle}*/}
                {/*            </strong>*/}
                {/*            <Progress*/}
                {/*                className="progress-sm"*/}
                {/*                value={userDetails.performanceReportValue}*/}
                {/*            >*/}
                {/*    <span className="progress-value">*/}
                {/*      {userDetails.performanceReportValue}%*/}
                {/*    </span>*/}
                {/*            </Progress>*/}
                {/*        </div>*/}
                {/*    </ListGroupItem>*/}
                {/*    <ListGroupItem className="p-4">*/}
                {/*        <strong className="text-muted d-block mb-2">*/}
                {/*            {userDetails.metaTitle}*/}
                {/*        </strong>*/}
                {/*        <span>{userDetails.metaValue}</span>*/}
                {/*    </ListGroupItem>*/}
                {/*</ListGroup>*/}
            </Card>
        )
    }


}

UserDetails.propTypes = {
    /**
     * The user details object.
     */
    data: PropTypes.object
};

UserDetails.defaultProps = {
    data: {
        name: "",
        avatar: require("./../../images/avatars/0.png"),
        jobTitle: "",
        performanceReportTitle: "",
        performanceReportValue: 0,
        // metaTitle: "Description",
        // metaValue:
        //     "Lorem ipsum dolor sit amet consectetur adipisicing elit. Odio eaque, quidem, commodi soluta qui quae minima obcaecati quod dolorum sint alias, possimus illum assumenda eligendi cumque?"
    }
};

export default UserDetails;
