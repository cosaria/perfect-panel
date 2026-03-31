'use client';

import { UserSubscribeDetail } from '@/app/dashboard/user/user-detail';
import { queryServerTotalData, queryTicketWaitReply } from '@/services/admin/console';
import { useQuery } from '@tanstack/react-query';
import { Card, CardContent, CardHeader, CardTitle } from '@workspace/ui/components/card';
import { ChartContainer, ChartTooltip, ChartTooltipContent } from '@workspace/ui/components/chart';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@workspace/ui/components/select';
import { Separator } from '@workspace/ui/components/separator';
import { Tabs, TabsList, TabsTrigger } from '@workspace/ui/components/tabs';
import { Icon } from '@workspace/ui/custom-components/icon';
import { formatBytes } from '@workspace/ui/utils';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { useState } from 'react';
import { Bar, BarChart, CartesianGrid, LabelList, XAxis, YAxis } from 'recharts';
import { Empty } from '../empty';
import { RevenueStatisticsCard } from './revenue-statistics-card';
import SystemVersionCard from './system-version-card';
import { UserStatisticsCard } from './user-statistics-card';

export default function Statistics() {
  const t = useTranslations('index');

  const { data: TicketTotal, isLoading: ticketLoading } = useQuery({
    queryKey: ['queryTicketWaitReply'],
    queryFn: async () => {
      const { data } = await queryTicketWaitReply();
      return data.data?.count;
    },
  });
  const { data: ServerTotal, isLoading: serverLoading } = useQuery({
    queryKey: ['queryServerTotalData'],
    queryFn: async () => {
      const { data } = await queryServerTotalData();
      return data.data;
    },
  });

  const isLoading = ticketLoading || serverLoading;

  const [dataType, setDataType] = useState<string | 'nodes' | 'users'>('nodes');
  const [timeFrame, setTimeFrame] = useState<string | 'today' | 'yesterday'>('today');

  const trafficData = {
    nodes: {
      today:
        ServerTotal?.server_traffic_ranking_today?.map((item) => ({
          name: item.name,
          traffic: item.download + item.upload,
        })) || [],
      yesterday:
        ServerTotal?.server_traffic_ranking_yesterday?.map((item) => ({
          name: item.name,
          traffic: item.download + item.upload,
        })) || [],
    },
    users: {
      today:
        ServerTotal?.user_traffic_ranking_today?.map((item) => ({
          name: item.sid,
          traffic: item.download + item.upload,
        })) || [],
      yesterday:
        ServerTotal?.user_traffic_ranking_yesterday?.map((item) => ({
          name: item.sid,
          traffic: item.download + item.upload,
        })) || [],
    },
  };
  const currentData =
    trafficData[dataType as 'nodes' | 'users'][timeFrame as 'today' | 'yesterday'];

  return (
    <>
      <div className='grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3'>
        {[
          {
            title: t('onlineUsersCount'),
            value: ServerTotal?.online_users || 0,
            subtitle: t('currentlyOnline'),
            icon: 'uil:users-alt',
            href: '/dashboard/user',
            color: 'text-blue-600 dark:text-blue-400',
            iconBg: 'bg-blue-100 dark:bg-blue-900/30',
          },

          {
            title: t('todayTraffic'),
            value: formatBytes(
              (ServerTotal?.today_upload || 0) + (ServerTotal?.today_download || 0),
            ),
            subtitle: `↑${formatBytes(ServerTotal?.today_upload || 0)} ↓${formatBytes(ServerTotal?.today_download || 0)}`,
            icon: 'uil:exchange-alt',
            color: 'text-purple-600 dark:text-purple-400',
            iconBg: 'bg-purple-100 dark:bg-purple-900/30',
          },
          {
            title: t('monthTraffic'),
            value: formatBytes(
              (ServerTotal?.monthly_upload || 0) + (ServerTotal?.monthly_download || 0),
            ),
            subtitle: `↑${formatBytes(ServerTotal?.monthly_upload || 0)} ↓${formatBytes(ServerTotal?.monthly_download || 0)}`,
            icon: 'uil:cloud-data-connection',
            color: 'text-orange-600 dark:text-orange-400',
            iconBg: 'bg-orange-100 dark:bg-orange-900/30',
          },
          {
            title: t('totalServers'),
            value: (ServerTotal?.online_servers || 0) + (ServerTotal?.offline_servers || 0),
            subtitle: `${t('online')} ${ServerTotal?.online_servers || 0} ${t('offline')} ${ServerTotal?.offline_servers || 0}`,
            icon: 'uil:server-network',
            href: '/dashboard/servers',
            color: 'text-green-600 dark:text-green-400',
            iconBg: 'bg-green-100 dark:bg-green-900/30',
          },
          {
            title: t('pendingTickets'),
            value: TicketTotal || 0,
            subtitle: t('pending'),
            icon: 'uil:clipboard-notes',
            href: '/dashboard/ticket',
            color: 'text-red-600 dark:text-red-400',
            iconBg: 'bg-red-100 dark:bg-red-900/30',
          },
        ].map((item, index) => (
          <Link
            href={item.href || '#'}
            key={index}
            className={!item.href ? 'pointer-events-none' : ''}
          >
            <Card className={`group ${item.href ? 'cursor-pointer' : ''}`}>
              <CardContent className='p-6'>
                <div className='flex items-center justify-between'>
                  <div className='flex-1'>
                    <p className='text-muted-foreground mb-2 text-sm font-medium'>{item.title}</p>
                    <div className={`text-2xl font-bold ${item.color} mb-1`}>{item.value}</div>
                    <div className={`text-muted-foreground h-4 text-xs`}>{item.subtitle}</div>
                  </div>
                  <div
                    className={`rounded-full p-3 ${item.iconBg} transition-transform duration-300 group-hover:scale-110`}
                  >
                    <Icon icon={item.icon} className={`h-6 w-6 ${item.color}`} />
                  </div>
                </div>
              </CardContent>
            </Card>
          </Link>
        ))}
        <SystemVersionCard />
      </div>
      <div className='grid gap-3 md:grid-cols-2 lg:grid-cols-3'>
        <RevenueStatisticsCard />
        <UserStatisticsCard />
        <Card>
          <CardHeader className='flex !flex-row items-center justify-between'>
            <CardTitle>{t('trafficRank')}</CardTitle>
            <Tabs value={timeFrame} onValueChange={setTimeFrame}>
              <TabsList>
                <TabsTrigger value='today'>{t('today')}</TabsTrigger>
                <TabsTrigger value='yesterday'>{t('yesterday')}</TabsTrigger>
              </TabsList>
            </Tabs>
          </CardHeader>
          <CardContent className='h-80'>
            <div className='mb-6 flex items-center justify-between'>
              <h4 className='font-semibold'>
                {dataType === 'nodes' ? t('nodeTraffic') : t('userTraffic')}
              </h4>
              <Select onValueChange={setDataType} defaultValue='nodes'>
                <SelectTrigger className='w-28'>
                  <SelectValue placeholder={t('selectTypePlaceholder')} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value='nodes'>{t('nodes')}</SelectItem>
                  <SelectItem value='users'>{t('users')}</SelectItem>
                </SelectContent>
              </Select>
            </div>
            {currentData.length > 0 ? (
              <ChartContainer
                config={{
                  traffic: {
                    label: t('traffic'),
                    color: 'hsl(var(--primary))',
                  },
                  type: {
                    label: t('type'),
                    color: 'hsl(var(--muted))',
                  },
                  label: {
                    color: 'hsl(var(--foreground))',
                  },
                }}
                className='max-h-80'
              >
                <BarChart data={currentData} layout='vertical' height={400}>
                  <CartesianGrid strokeDasharray='3 3' />
                  <XAxis
                    type='number'
                    tickLine={false}
                    axisLine={false}
                    tickFormatter={(value) => formatBytes(value || 0)}
                  />
                  <YAxis
                    type='category'
                    dataKey='name'
                    tickLine={false}
                    axisLine={false}
                    interval={0}
                    tickMargin={0}
                    width={15}
                    tickFormatter={(value, index) => String(index + 1)}
                  />
                  <ChartTooltip
                    trigger='hover'
                    content={
                      <ChartTooltipContent
                        label={true}
                        labelFormatter={(label, [payload]) =>
                          dataType === 'nodes' ? (
                            `${t('nodes')}: ${label}`
                          ) : (
                            <>
                              <div className='w-80'>
                                <UserSubscribeDetail id={payload?.payload.name} enabled={true} />
                              </div>
                              <Separator className='my-2' />
                              <div>{`${t('users')}: ${label}`}</div>
                            </>
                          )
                        }
                        formatter={(value) => {
                          return formatBytes(Number(value) || 0);
                        }}
                      />
                    }
                  />
                  <Bar dataKey='traffic' fill='hsl(var(--primary))' radius={[0, 4, 4, 0]}>
                    <LabelList
                      dataKey='name'
                      position='insideLeft'
                      offset={8}
                      className='fill-[--color-label]'
                      fontSize={12}
                    />
                  </Bar>
                </BarChart>
              </ChartContainer>
            ) : (
              <div className='flex h-full items-center justify-center'>
                <Empty />
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </>
  );
}
