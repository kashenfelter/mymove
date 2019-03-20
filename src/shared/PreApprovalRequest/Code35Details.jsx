import React from 'react';
import { formatCents } from 'shared/formatters';

export const Code35Details = props => {
  const row = props.shipmentLineItem;
  const actAmtValue = formatCents(row.actual_amount_cents);
  const actAmtText = actAmtValue ? '$' + actAmtValue : '--';
  return (
    <td details-cy={`${row.tariff400ng_item.code}-details`}>
      {row.description} <br />
      {row.reason} <br />
      Est. not to exceed: ${formatCents(row.estimate_amount_cents)} <br />
      Actual amount: {actAmtText} <br />
      {row.notes}
    </td>
  );
};
