'use client'

import { AreaChart, Card, Title } from '@tremor/react'

const generateData = () => {
  const dataset = []
  const dates = [
    'Jun 30',
    'Jul 01',
    'Jul 02',
    'Jul 03',
    'Jul 04',
    'Jul 05',
    'Jul 06',
    'Jul 07',
    'Jul 08',
    'Jul 09',
    'Jul 10',
    'Jul 11',
    'Jul 12',
    'Jul 13',
    'Jul 14',
    'Jul 15',
    'Jul 16',
    'Jul 17',
  ]

  for (const date of dates) {
    dataset.push({
      date,
      'checkout-1': Math.round(150 + Math.random() * 20 - 10),
      'checkout-2': Math.round(200 + Math.random() * 20 - 10),
      'checkout-3': Math.round(250 + Math.random() * 20 - 10),
    })
  }

  return dataset
}

const mockDataset = generateData()
console.log(mockDataset)

export default function Chart() {
  return (
    <Card className='mt-8'>
      <Title className='mb-2'>My admin dashboard</Title>
      <AreaChart
        className='mt-4 h-80'
        defaultValue={0}
        data={mockDataset}
        categories={['checkout-1', 'checkout-2', 'checkout-3']}
        index='date'
        colors={['indigo', 'fuchsia', 'emerald', 'neutral']}
        allowDecimals={false}
        yAxisWidth={60}
        noDataText='No data. Run your first test to get started!'
      />
    </Card>
  )
}