import React from 'react';
import { shallow } from 'enzyme';
import PreApprovalTable from './PreApprovalTable';

describe('PreApprovalTable tests', () => {
  let wrapper, icons;
  const onEdit = jest.fn();
  const shipmentAccessorials = [
    {
      id: 'sldkjf',
      accessorial: { code: '105D', item: 'Reg Shipping' },
      location: 'D',
      base_quantity: 167000,
      notes: '',
      created_at: '2018-09-24T14:05:38.847Z',
      status: 'SUBMITTED',
    },
    {
      id: 'sldsdff',
      accessorial: { code: '105D', item: 'Reg Shipping' },
      location: 'D',
      base_quantity: 788300,
      notes: 'Mounted deer head measures 23" x 34" x 27"; crate will be 16.7 cu ft',
      created_at: '2018-09-24T14:05:38.847Z',
      status: 'SUBMITTED',
    },
  ];
  describe('When on approval is passed in and status is submitted', () => {
    it('renders without crashing', () => {
      wrapper = shallow(
        <PreApprovalTable
          shipmentAccessorials={shipmentAccessorials}
          isActionable={true}
          onEdit={onEdit}
          onDelete={onEdit}
          onApproval={onEdit}
        />,
      );
      const icons = wrapper.find('.icon');
      expect(wrapper.find('.accessorial-panel').length).toEqual(1);
      expect(icons.length).toBe(2);
    });
  });
  describe('When on approval is NOT passed in and status is SUBMITTED', () => {
    beforeEach(() => {
      wrapper = shallow(
        <PreApprovalTable
          shipmentAccessorials={shipmentAccessorials}
          isActionable={true}
          onEdit={onEdit}
          onDelete={onEdit}
        />,
      );
    });
    it('it shows the appropriate number of icons.', () => {
      const icons = wrapper.find('.icon');
      expect(icons.length).toBe(0);
    });
  });
  describe('When on approval is passed in and status is APPROVED', () => {
    beforeEach(() => {
      shipmentAccessorials[0].status = 'APPROVED';
      shipmentAccessorials[1].status = 'APPROVED';
      wrapper = shallow(
        <PreApprovalTable
          shipmentAccessorials={shipmentAccessorials}
          isActionable={true}
          onEdit={onEdit}
          onDelete={onEdit}
          onApproval={onEdit}
        />,
      );
    });
    it('it shows the appropriate number of icons.', () => {
      const icons = wrapper.find('.icon');
      expect(icons.length).toBe(0);
    });
  });
  describe('When on approval is NOT passed in and status is APPROVED', () => {
    beforeEach(() => {
      shipmentAccessorials[0].status = 'APPROVED';
      shipmentAccessorials[1].status = 'APPROVED';
      wrapper = shallow(
        <PreApprovalTable
          shipmentAccessorials={shipmentAccessorials}
          isActionable={true}
          onEdit={onEdit}
          onDelete={onEdit}
        />,
      );
    });
    it('it shows the appropriate number of icons.', () => {
      const icons = wrapper.find('.icon');
      expect(icons.length).toBe(0);
    });
  });
});
