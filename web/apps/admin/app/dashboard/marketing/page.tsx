'use client';

import { Table, TableBody, TableCell, TableRow } from '@workspace/ui/components/table';
import { useTranslations } from 'next-intl';
import EmailBroadcastForm from './email/broadcast-form';
import EmailTaskManager from './email/task-manager';
import QuotaBroadcastForm from './quota/broadcast-form';
import QuotaTaskManager from './quota/task-manager';

export default function MarketingPage() {
  const t = useTranslations('marketing');

  const formSections = [
    {
      title: t('emailMarketing'),
      forms: [{ component: EmailBroadcastForm }, { component: EmailTaskManager }],
    },
    {
      title: t('quotaService'),
      forms: [{ component: QuotaBroadcastForm }, { component: QuotaTaskManager }],
    },
  ];

  return (
    <div className='space-y-8'>
      {formSections.map((section, sectionIndex) => (
        <div key={sectionIndex}>
          <h2 className='mb-4 text-lg font-semibold'>{section.title}</h2>
          <Table>
            <TableBody>
              {section.forms.map((form, formIndex) => {
                const FormComponent = form.component;
                return (
                  <TableRow key={formIndex}>
                    <TableCell>
                      <FormComponent />
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </div>
      ))}
    </div>
  );
}
