import React from "react";
import PropTypes from "prop-types";
import {
    Card,
    CardHeader,
    ListGroup,
    ListGroupItem,
    Row,
    Col,
    Form,
    FormGroup,
    FormInput,
    FormSelect,
    FormTextarea,
    Button
} from "shards-react";


const UserAccountDetails = ({title, data}) => (
    <Card small className="mb-4">
        <CardHeader className="border-bottom">
            <h6 className="m-0">{title}</h6>
        </CardHeader>
        <ListGroup flush>
            <ListGroupItem className="p-3">
                <Row>
                    <Col>
                        <Form>
                            <Row form>
                                {/* First Name */}
                                <Col md="6" className="form-group">
                                    <label htmlFor="feFirstName">First Name</label>
                                    <FormInput
                                        id="feFirstName"
                                        placeholder="First Name"
                                        value={data.first_name}
                                        onChange={() => {
                                        }}
                                    />
                                </Col>
                                {/* Last Name */}
                                <Col md="6" className="form-group">
                                    <label htmlFor="feLastName">Last Name</label>
                                    <FormInput
                                        id="feLastName"
                                        placeholder="Last Name"
                                        value={data.last_name}
                                        onChange={() => {
                                        }}
                                    />
                                </Col>
                            </Row>
                            <Row form>
                                {/* Email */}
                                <Col md="12" className="form-group">
                                    <label htmlFor="feEmail">Email</label>
                                    <FormInput
                                        type="email"
                                        id="feEmail"
                                        placeholder="Email Address"
                                        value={data.email}
                                        onChange={() => {
                                        }}
                                        autoComplete="email"
                                    />
                                </Col>
                                {/* Password */}
                                {/*<Col md="6" className="form-group">*/}
                                {/*    <label htmlFor="fePassword">Password</label>*/}
                                {/*    <FormInput*/}
                                {/*        type="password"*/}
                                {/*        id="fePassword"*/}
                                {/*        placeholder="Password"*/}
                                {/*        value="EX@MPL#P@$$w0RD"*/}
                                {/*        onChange={() => {*/}
                                {/*        }}*/}
                                {/*        autoComplete="current-password"*/}
                                {/*    />*/}
                                {/*</Col>*/}
                            </Row>
                            <FormGroup>
                                <label htmlFor="feAddress">Position</label>
                                <FormInput
                                    id="feAddress"
                                    placeholder="Your position or title"
                                    value={data.position}
                                    onChange={() => {
                                    }}
                                />
                            </FormGroup>
                            <FormGroup>
                                <label htmlFor="feAddress">Company Name</label>
                                <FormInput
                                    id="compName"
                                    placeholder="Company Name"
                                    value={data.company_name}
                                    onChange={() => {
                                    }}
                                />
                            </FormGroup>
                            <Row form>
                                {/* City */}
                                <Col md="6" className="form-group">
                                    <label htmlFor="feCity">City</label>
                                    <FormInput
                                        id="feCity"
                                        placeholder="City"
                                        value={data.city}
                                        onChange={() => {
                                        }}
                                    />
                                </Col>
                                {/* State */}
                                <Col md="6" className="form-group">
                                    <label htmlFor="feCountry">Country</label>
                                    <FormInput
                                        id="feCountry"
                                        placeholder="Country"
                                        value={data.country}
                                        onChange={() => {
                                        }}
                                    />
                                </Col>
                                {/* Zip Code */}
                            </Row>
                            {/*<Row form>*/}
                            {/*    /!* Description *!/*/}
                            {/*    <Col md="12" className="form-group">*/}
                            {/*        <label htmlFor="feDescription">Description</label>*/}
                            {/*        <FormTextarea id="feDescription" rows="5"/>*/}
                            {/*    </Col>*/}
                            {/*</Row>*/}
                            <Button disabled theme="accent">Update Account</Button>
                        </Form>
                    </Col>
                </Row>
            </ListGroupItem>
        </ListGroup>
    </Card>
);

UserAccountDetails.propTypes = {
    /**
     * The component's title.
     */
    title: PropTypes.string
};

UserAccountDetails.defaultProps = {
    title: "Account Details"
};

export default UserAccountDetails;
