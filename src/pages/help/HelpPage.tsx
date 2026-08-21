import { useState } from 'react';
import { SearchIcon, ChevronDownIcon, ChevronUpIcon, LifeBuoyIcon, BookOpenIcon, MessageSquareIcon, TicketIcon, ArrowRightIcon, CheckCircleIcon, ClockIcon, UserIcon, ShieldIcon, ZapIcon, GlobeIcon, MailIcon, PhoneIcon, VideoIcon } from 'lucide-react';
import { cn } from '../../utils/helpers';
import { Button } from '../../components/ui/Button';
import { Input } from '../../components/ui/Input';
import { Card, CardHeader, CardBody, CardTitle, CardDescription } from '../../components/ui/Card';
import { Badge } from '../../components/ui/Badge';
import { Tabs, TabItem } from '../../components/ui/Tabs';
import { Breadcrumb } from '../../components/ui/Breadcrumb';
import { Accordion, AccordionItem, AccordionTrigger, AccordionContent } from '../../components/ui/Accordion';

const faqs = [
  {
    category: 'Getting Started',
    items: [
      { question: 'How do I set up my organization?', answer: 'After signing up, you\'ll be guided through the onboarding wizard. Enter your organization details, add facilities, and connect data sources. The process takes about 15 minutes.' },
      { question: 'What data do I need to get started?', answer: 'At minimum, you\'ll need facility locations and energy bills. For comprehensive tracking, gather utility data, fuel consumption records, and operational metrics.' },
      { question: 'Can I import historical data?', answer: 'Yes, you can import up to 10 years of historical data via CSV, Excel, or API. Our data validation engine will flag any inconsistencies.' },
    ],
  },
  {
    category: 'Emissions Tracking',
    items: [
      { question: 'Which GHG Protocol standards are supported?', answer: 'We support Corporate Standard, Scope 2 Guidance, Scope 3 Standard, and Product Standard. All calculations follow the latest GHG Protocol methodologies.' },
      { question: 'How are emission factors updated?', answer: 'Emission factors are automatically updated from EPA eGRID, IEA, DEFRA, and IPCC databases. You can also override with custom factors.' },
      { question: 'Can I track Scope 3 categories?', answer: 'Yes, all 15 Scope 3 categories are supported with category-specific calculation methodologies and supplier engagement tools.' },
    ],
  },
  {
    category: 'Integrations & Data',
    items: [
      { question: 'Which utility providers are supported?', answer: 'We have direct integrations with 200+ utility providers across North America, Europe, and APAC. Custom integrations are available via API.' },
      { question: 'How does automated data collection work?', answer: 'Connect your utility accounts, IoT sensors, or ERP systems. Data is pulled daily, validated, and normalized automatically.' },
      { question: 'What file formats are supported for upload?', answer: 'CSV, Excel (.xlsx, .xls), JSON, and PDF (with OCR for invoices). Maximum file size is 100MB.' },
    ],
  },
  {
    category: 'Reporting & Compliance',
    items: [
      { question: 'Which reporting frameworks are supported?', answer: 'CDP, TCFD, GRI, SASB, SFDR, CSRD, SEC Climate Rules, ISSB, and custom frameworks. Reports auto-populate from your data.' },
      { question: 'Can I customize report templates?', answer: 'Yes, use our drag-and-drop report builder or modify existing templates. Branding, charts, and narrative sections are fully customizable.' },
      { question: 'How do I prepare for assurance?', answer: 'Our audit trail tracks all data changes with timestamps and user IDs. Export the complete audit log for verifiers.' },
    ],
  },
  {
    category: 'Account & Billing',
    items: [
      { question: 'What payment methods are accepted?', answer: 'Credit cards (Visa, Mastercard, Amex), ACH/wire transfers for annual plans, and purchase orders for enterprise.' },
      { question: 'Can I change my plan mid-cycle?', answer: 'Upgrades take effect immediately with prorated billing. Downgrades take effect at the next billing cycle.' },
      { question: 'Is there a free trial?', answer: 'Yes, 14-day full-featured trial with no credit card required. All features including integrations are available.' },
    ],
  },
];

const contactOptions = [
  { icon: MessageSquareIcon, title: 'Live Chat', description: 'Available 9am-6pm EST, Mon-Fri', action: 'Start Chat', color: 'text-green-500' },
  { icon: MailIcon, title: 'Email Support', description: 'Response within 4 hours', action: 'Email Us', color: 'text-blue-500', href: 'mailto:support@greencombats.io' },
  { icon: TicketIcon, title: 'Submit Ticket', description: 'For technical issues and feature requests', action: 'Create Ticket', color: 'text-purple-500' },
  { icon: VideoIcon, title: 'Schedule Demo', description: '30-min personalized walkthrough', action: 'Book Demo', color: 'text-primary-500' },
  { icon: PhoneIcon, title: 'Phone Support', description: 'Enterprise plans only', action: 'Call Us', color: 'text-orange-500' },
  { icon: LifeBuoyIcon, title: 'Help Center', description: 'Search 200+ articles and guides', action: 'Browse Docs', color: 'text-secondary-500', href: '/docs' },
];

const resources = [
  { icon: BookOpenIcon, title: 'Documentation', description: 'Complete API reference and user guides', link: '/docs' },
  { icon: VideoIcon, title: 'Video Tutorials', description: 'Step-by-step walkthroughs for every feature', link: '/academy' },
  { icon: ZapIcon, title: 'API Reference', description: 'REST and GraphQL API documentation', link: '/api-docs' },
  { icon: GlobeIcon, title: 'Integration Guides', description: 'Setup guides for all 300+ integrations', link: '/integrations/guides' },
  { icon: ShieldIcon, title: 'Compliance Templates', description: 'Ready-to-use templates for major frameworks', link: '/templates' },
  { icon: UserIcon, title: 'Community Forum', description: 'Connect with other climate professionals', link: '/community' },
];

export default function HelpPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [activeCategory, setActiveCategory] = useState('all');

  const allFAQs = faqs.flatMap((cat) => cat.items.map((item) => ({ ...item, category: cat.category })));
  const filteredFAQs = allFAQs.filter((faq) => 
    (activeCategory === 'all' || faq.category === activeCategory) &&
    (faq.question.toLowerCase().includes(searchQuery.toLowerCase()) || 
     faq.answer.toLowerCase().includes(searchQuery.toLowerCase()))
  );

  return (
    <div className="space-y-6">
      <div>
        <Breadcrumb items={[
          { label: 'Help & Support', href: '/help', current: true },
        ]} />
        <h1 className="mt-2 text-2xl font-bold text-secondary-900 dark:text-white">Help & Support</h1>
        <p className="text-secondary-600 dark:text-secondary-400">Find answers, explore resources, or contact our team</p>
      </div>

      {/* Search & Quick Actions */}
      <Card>
        <CardBody className="p-6">
          <div className="flex flex-col lg:flex-row gap-6 items-start lg:items-center justify-between">
            <div className="relative max-w-xl w-full">
              <SearchIcon className="absolute left-4 top-1/2 h-5 w-5 -translate-y-1/2 text-secondary-400" />
              <Input
                placeholder="Search help articles, FAQs, and guides..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-12"
              />
            </div>
            <div className="flex flex-wrap gap-3">
              {['all', 'Getting Started', 'Emissions Tracking', 'Integrations & Data', 'Reporting & Compliance', 'Account & Billing'].map((cat) => (
                <button
                  key={cat}
                  onClick={() => setActiveCategory(cat)}
                  className={cn(
                    'px-4 py-2 rounded-lg text-sm font-medium transition-colors',
                    activeCategory === cat
                      ? 'bg-primary-100 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300'
                      : 'bg-secondary-100 text-secondary-700 hover:bg-secondary-200 dark:bg-secondary-800 dark:text-secondary-300 dark:hover:bg-secondary-700'
                  )}
                >
                  {cat === 'all' ? 'All Topics' : cat}
                </button>
              ))}
            </div>
          </div>
        </CardBody>
      </Card>

      {/* Contact Options */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {contactOptions.map((option, i) => (
          <Card key={i} className="hover:shadow-md transition-shadow h-full">
            <CardBody className="p-6">
              <div className="flex items-start gap-4">
                <div className={cn('p-3 rounded-xl', option.color.replace('text-', 'bg-') + '/10')}>
                  <option.icon className={cn('h-6 w-6', option.color)} />
                </div>
                <div className="flex-1">
                  <h3 className="font-semibold text-secondary-900 dark:text-white">{option.title}</h3>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400 mt-1">{option.description}</p>
                  <Button variant="outline" size="sm" className="mt-3" asChild>
                    {option.href ? <a href={option.href}>{option.action}</a> : <button>{option.action}</button>}
                  </Button>
                </div>
              </div>
            </CardBody>
          </Card>
        ))}
      </div>

      {/* FAQ Sections */}
      <div className="space-y-8">
        {faqs.map((category) => (
          <Card key={category.category}>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle>{category.category}</CardTitle>
                  <CardDescription>{category.items.length} frequently asked questions</CardDescription>
                </div>
                <Badge variant="outline" size="sm">{category.items.length} articles</Badge>
              </div>
            </CardHeader>
            <CardBody className="pt-0">
              <Accordion type="single" collapsible className="w-full">
                {category.items.map((faq, i) => (
                  <AccordionItem key={i} value={faq.question}>
                    <AccordionTrigger className="text-left py-4">
                      {faq.question}
                    </AccordionTrigger>
                    <AccordionContent className="pb-4 text-secondary-600 dark:text-secondary-400">
                      {faq.answer}
                    </AccordionContent>
                  </AccordionItem>
                ))}
              </Accordion>
            </CardBody>
          </Card>
        ))}
      </div>

      {/* Resources */}
      <Card>
        <CardHeader>
          <CardTitle>Resources & Learning</CardTitle>
          <CardDescription>Deepen your knowledge with our comprehensive resource library</CardDescription>
        </CardHeader>
        <CardBody>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {resources.map((resource, i) => (
              <Link key={i} to={resource.link} className="flex items-center gap-4 p-4 rounded-lg border border-secondary-200 dark:border-secondary-700 hover:border-primary-300 dark:hover:border-primary-700 transition-colors group">
                <div className={cn('p-3 rounded-xl', 'bg-primary-100 dark:bg-primary-900/30')}>
                  <resource.icon className="h-6 w-6 text-primary-600 dark:text-primary-400 group-hover:scale-110 transition-transform" />
                </div>
                <div>
                  <p className="font-medium text-secondary-900 dark:text-white group-hover:text-primary-600 dark:group-hover:text-primary-400 transition-colors">{resource.title}</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">{resource.description}</p>
                </div>
              </Link>
            ))}
          </div>
        </CardBody>
      </Card>

      {/* Status & SLA */}
      <Card className="bg-primary-50 dark:bg-primary-900/20 border-primary-200 dark:border-primary-800">
        <CardBody className="p-6">
          <div className="flex items-center gap-4">
            <div className="p-3 bg-primary-100 dark:bg-primary-900/30 rounded-xl">
              <CheckCircleIcon className="h-6 w-6 text-primary-600 dark:text-primary-400" />
            </div>
            <div className="flex-1">
              <h3 className="font-semibold text-primary-900 dark:text-primary-100">System Status: All Systems Operational</h3>
              <p className="text-sm text-primary-700 dark:text-primary-300 mt-1">Last updated: Just now</p>
            </div>
            <Button variant="secondary" size="sm" asChild>
              <a href="/status">View Status Page <ArrowRightIcon className="h-4 w-4 ml-2" /></a>
            </Button>
          </div>
        </CardBody>
      </Card>
    </div>
  );
}