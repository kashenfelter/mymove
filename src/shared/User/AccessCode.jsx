import React, { Fragment } from 'react';
import { bindActionCreators } from 'redux';
import { getFormValues, reduxForm, SubmissionError } from 'redux-form';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { get } from 'lodash';
import PropTypes from 'prop-types';
import './AccessCode.css';
import { validateAccessCode, claimAccessCode } from 'shared/Entities/modules/accessCodes';

class AccessCode extends React.Component {
  validateAndClaimAccessCode = () => {
    const { formValues, validateAccessCode, claimAccessCode } = this.props;
    return validateAccessCode(formValues.claim_access_code)
      .then(res => {
        const { body: accessCode } = get(res, 'response');
        if (!accessCode.code) {
          throw new SubmissionError({
            claim_access_code: 'This code is invalid',
            _error: 'Validating access code failed!',
          });
        }
        claimAccessCode(accessCode)
          .then(() => {
            window.location.reload();
          })
          .catch(err => {
            console.log(err);
          });
      })
      .catch(err => {
        throw err;
      });
  };
  render() {
    const { schema } = this.props;
    return (
      <Fragment>
        <div className="usa-grid">
          <h3 className="title">Welcome to MilMove</h3>
          <p>Please enter your MilMove access code in the field below.</p>
          <SwaggerField fieldName="claim_access_code" swagger={schema} required />
          <button className="usa-button-primary" onClick={this.validateAndClaimAccessCode}>
            Continue
          </button>
          <br />No code? Go to DPS to schedule your move.
        </div>
      </Fragment>
    );
  }
}

const formName = 'claim_access_code_form';
AccessCode = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(AccessCode);

AccessCode.propTypes = {
  schema: PropTypes.object.isRequired,
  serviceMemberId: PropTypes.string.isRequired,
};

function mapStateToProps(state) {
  const serviceMember = get(state, 'serviceMember.currentServiceMember');
  const props = {
    schema: get(state, 'swaggerInternal.spec.definitions.ClaimAccessCodePayload', {}),
    //accessCode: get(state, 'entities.validateAccessCode.undefined.access_code', {}),
    serviceMemberId: get(serviceMember, 'id'),
    formValues: getFormValues(formName)(state),
  };
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ validateAccessCode, claimAccessCode }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(withRouter(AccessCode));
