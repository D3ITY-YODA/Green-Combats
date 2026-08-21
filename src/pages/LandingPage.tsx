import { Link } from 'react-router-dom';
import { ArrowRightIcon, CheckIcon, BarChart2Icon, MapIcon, LeafIcon, ZapIcon, DropletIcon, CloudIcon, ShieldIcon, UsersIcon, GlobeIcon, AwardIcon, ChevronRightIcon } from 'lucide-react';
import { Button } from '../components/ui/Button';
import { Card, CardBody } from '../components/ui/Card';
import { Badge } from '../components/ui/Badge';

const features = [
  { icon: BarChart2Icon, title: 'Real-time Analytics', description: 'Track emissions, energy, water, and biodiversity metrics with live dashboards and automated reporting.' },
  { icon: MapIcon, title: 'Interactive Maps', description: 'Visualize climate data on geospatial maps with satellite imagery, facility locations, and risk zones.' },
  { icon: LeafIcon, title: 'Carbon Accounting', description: 'Complete GHG Protocol compliant Scope 1, 2, and 3 emissions tracking with automated calculations.' },
  { icon: ZapIcon, title: 'Energy Management', description: 'Monitor renewable energy adoption, efficiency gains, and grid interactions across all facilities.' },
  { icon: DropletIcon, title: 'Water Stewardship', description: 'Track water withdrawal, consumption, and discharge with watershed-level risk assessment.' },
  { icon: CloudIcon, title: 'Air Quality Monitoring', description: 'Real-time air quality indices with pollutant tracking and health impact assessments.' },
  { icon: ShieldIcon, title: 'Regulatory Compliance', description: 'Automated compliance tracking for CSRD, TCFD, CDP, SEC, and regional regulations.' },
  { icon: UsersIcon, title: 'Collaboration Tools', description: 'Team workspaces, task management, and stakeholder engagement for climate initiatives.' },
  { icon: GlobeIcon, title: 'Scenario Modeling', description: 'Climate scenario analysis with IPCC-aligned pathways and financial impact projections.' },
  { icon: AwardIcon, title: 'Certification Ready', description: 'Built-in templates for Science Based Targets, RE100, and net-zero certifications.' },
];

const stats = [
  { value: '500+', label: 'Organizations' },
  { value: '50K+', label: 'Facilities Tracked' },
  { value: '12M+', label: 'tCO₂e Managed' },
  { value: '99.9%', label: 'Uptime SLA' },
];

const integrations = [
  { name: 'Utility Providers', count: 200 },
  { name: 'IoT Sensors', count: 50 },
  { name: 'ERP Systems', count: 15 },
  { name: 'Satellite Data', count: 10 },
  { name: 'Weather APIs', count: 8 },
  { name: 'Carbon Registries', count: 12 },
];

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-white dark:bg-secondary-950">
      {/* Navigation Bar */}
      <header className="fixed top-0 left-0 right-0 z-50 bg-white/80 dark:bg-secondary-950/80 backdrop-blur-lg border-b border-secondary-200 dark:border-secondary-800">
        <nav className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8" aria-label="Main navigation">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center gap-8">
              <Link to="/" className="flex items-center gap-2" aria-label="Green Combats Home">
                <div className="h-10 w-10 rounded-xl bg-primary-600 flex items-center justify-center">
                  <span className="text-white font-bold text-xl">GC</span>
                </div>
                <span className="text-xl font-bold text-secondary-900 dark:text-white hidden sm:block">Green Combats</span>
              </Link>
              <div className="hidden md:flex items-center gap-6">
                <Link to="#features" className="text-sm font-medium text-secondary-600 hover:text-secondary-900 dark:text-secondary-400 dark:hover:text-white transition-colors">Features</Link>
                <Link to="#solutions" className="text-sm font-medium text-secondary-600 hover:text-secondary-900 dark:text-secondary-400 dark:hover:text-white transition-colors">Solutions</Link>
                <Link to="#integrations" className="text-sm font-medium text-secondary-600 hover:text-secondary-900 dark:text-secondary-400 dark:hover:text-white transition-colors">Integrations</Link>
                <Link to="#pricing" className="text-sm font-medium text-secondary-600 hover:text-secondary-900 dark:text-secondary-400 dark:hover:text-white transition-colors">Pricing</Link>
                <Link to="#resources" className="text-sm font-medium text-secondary-600 hover:text-secondary-900 dark:text-secondary-400 dark:hover:text-white transition-colors">Resources</Link>
              </div>
            </div>
            <div className="flex items-center gap-4">
              <Link to="/auth/login" className="hidden sm:block text-sm font-medium text-secondary-600 hover:text-secondary-900 dark:text-secondary-400 dark:hover:text-white transition-colors">Sign in</Link>
              <Link to="/auth/register" className="hidden sm:block">
                <Button variant="primary" size="sm">Get Started Free</Button>
              </Link>
            </div>
          </div>
        </nav>
      </header>

      {/* Hero Section */}
      <section className="relative pt-32 pb-20 px-4 sm:px-6 lg:px-8 overflow-hidden">
        <div className="max-w-7xl mx-auto">
          <div className="grid lg:grid-cols-2 gap-12 items-center">
            <div className="text-center lg:text-left">
              <Badge variant="primary" className="mb-6 inline-flex" style={{ backgroundColor: 'rgba(34, 197, 94, 0.1)', color: '#16a34a' }}>
                <span className="mr-2">New</span>
                AI-Powered Climate Intelligence v2.0 Released
              </Badge>
              <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold text-secondary-900 dark:text-white leading-tight mb-6">
                Climate Action Intelligence
                <br />
                <span className="text-primary-600 dark:text-primary-400">That Drives Results</span>
              </h1>
              <p className="text-lg sm:text-xl text-secondary-600 dark:text-secondary-300 mb-8 max-w-2xl mx-auto lg:mx-0">
                The only platform that unifies emissions tracking, energy management, water stewardship, 
                and regulatory compliance into a single intelligent system. Trusted by 500+ organizations worldwide.
              </p>
              <div className="flex flex-col sm:flex-row items-center justify-center lg:justify-start gap-4 mb-12">
                <Link to="/auth/register">
                  <Button variant="primary" size="lg" className="w-full sm:w-auto gap-2">
                    Start Free Trial
                    <ArrowRightIcon className="h-5 w-5" />
                  </Button>
                </Link>
                <Link to="#demo">
                  <Button variant="outline" size="lg" className="w-full sm:w-auto gap-2">
                    Watch Demo
                    <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 20 20"><path d="M6.3 2.841A1.5 1.5 0 004 4.11V15.89a1.5 1.5 0 002.3 1.269l9.344-5.89a1.5 1.5 0 000-2.538L6.3 2.84z" /></svg>
                  </Button>
                </Link>
              </div>
              <p className="text-sm text-secondary-500 dark:text-secondary-400">
                No credit card required • 14-day free trial • Cancel anytime
              </p>
            </div>

            <div className="relative">
              <div className="bg-secondary-900 dark:bg-secondary-950 rounded-2xl border border-secondary-700 dark:border-secondary-800 overflow-hidden shadow-2xl">
                <div className="flex items-center gap-2 px-4 py-3 border-b border-secondary-700">
                  <div className="flex gap-1.5">
                    <div className="w-3 h-3 rounded-full bg-red-500" />
                    <div className="w-3 h-3 rounded-full bg-yellow-500" />
                    <div className="w-3 h-3 rounded-full bg-green-500" />
                  </div>
                  <div className="flex-1 text-center text-sm text-secondary-500 font-mono">app.greencombats.io/dashboard</div>
                </div>
                <div className="p-6 h-96 overflow-hidden">
                  <div className="grid grid-cols-3 gap-4 mb-6">
                    {[
                      { label: 'Total Emissions', value: '12,450', unit: 'tCO₂e', trend: -2.3, color: 'text-red-400' },
                      { label: 'Renewable Energy', value: '67%', unit: '', trend: 12.5, color: 'text-green-400' },
                      { label: 'Compliance Score', value: '94%', unit: '', trend: 2.1, color: 'text-blue-400' },
                    ].map((stat, i) => (
                      <div key={i} className="bg-secondary-800/50 rounded-xl p-4 border border-secondary-700">
                        <p className="text-sm text-secondary-400 mb-1">{stat.label}</p>
                        <p className="text-2xl font-bold text-white">{stat.value}<span className="text-lg font-normal">{stat.unit}</span></p>
                        <p className={cn('text-sm font-medium', stat.color)}>{stat.trend >= 0 ? '+' : ''}{stat.trend}%</p>
                      </div>
                    ))}
                  </div>
                  <div className="bg-secondary-800/50 rounded-xl p-4 border border-secondary-700 h-48 flex items-end justify-around">
                    {[
                      4200, 4100, 4000, 3900, 3850, 3780, 3720, 3650, 3580, 3520, 3450, 3400
                    ].map((val, i) => (
                      <div key={i} className="flex-1 max-w-8">
                        <div className="h-full bg-primary-500/50 rounded-t transition-all hover:bg-primary-500" style={{ height: `${(val / 4200) * 100}%` }} />
                        <span className="text-xs text-secondary-500 mt-1 block text-center">{['J','F','M','A','M','J','J','A','S','O','N','D'][i]}</span>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
              <div className="absolute -bottom-6 -right-6 lg:-right-12 w-72 h-72 bg-primary-500/20 rounded-full blur-3xl" aria-hidden="true" />
              <div className="absolute -top-6 -left-6 lg:-left-12 w-72 h-72 bg-blue-500/20 rounded-full blur-3xl" aria-hidden="true" />
            </div>
          </div>
        </div>
      </section>

      {/* Stats Bar */}
      <section className="py-16 px-4 sm:px-6 lg:px-8 bg-secondary-50 dark:bg-secondary-900/50 border-y border-secondary-200 dark:border-secondary-800">
        <div className="max-w-7xl mx-auto">
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-8">
            {stats.map((stat, i) => (
              <div key={i} className="text-center">
                <p className="text-3xl sm:text-4xl font-bold text-secondary-900 dark:text-white">{stat.value}</p>
                <p className="text-secondary-600 dark:text-secondary-400 mt-1">{stat.label}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="py-20 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto">
          <div className="text-center mb-16">
            <Badge variant="outline" className="mb-4">Features</Badge>
            <h2 className="text-3xl sm:text-4xl font-bold text-secondary-900 dark:text-white mb-4">
              Everything You Need for Climate Action
            </h2>
            <p className="text-lg text-secondary-600 dark:text-secondary-300 max-w-2xl mx-auto">
              Comprehensive tools for measuring, managing, and reducing your environmental impact.
            </p>
          </div>
          <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {features.map((feature, i) => (
              <Card key={i} className="hover:shadow-xl transition-shadow h-full">
                <CardBody className="p-6">
                  <div className="p-3 bg-primary-100 dark:bg-primary-900/30 rounded-xl w-fit mb-4">
                    <feature.icon className="h-6 w-6 text-primary-600 dark:text-primary-400" />
                  </div>
                  <h3 className="text-xl font-semibold text-secondary-900 dark:text-white mb-2">{feature.title}</h3>
                  <p className="text-secondary-600 dark:text-secondary-400">{feature.description}</p>
                </CardBody>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* Solutions Section */}
      <section id="solutions" className="py-20 px-4 sm:px-6 lg:px-8 bg-secondary-50 dark:bg-secondary-900/50">
        <div className="max-w-7xl mx-auto">
          <div className="text-center mb-16">
            <Badge variant="outline" className="mb-4">Solutions</Badge>
            <h2 className="text-3xl sm:text-4xl font-bold text-secondary-900 dark:text-white mb-4">
              Tailored for Your Industry
            </h2>
            <p className="text-lg text-secondary-600 dark:text-secondary-300 max-w-2xl mx-auto">
              Industry-specific configurations and templates for accelerated implementation.
            </p>
          </div>
          <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {[
              { title: 'Manufacturing', metrics: ['Scope 1&2 Tracking', 'Energy Optimization', 'Waste Reduction'], icon: ZapIcon },
              { title: 'Real Estate', metrics: ['Building Emissions', 'Tenant Engagement', 'GRESB Reporting'], icon: LeafIcon },
              { title: 'Financial Services', metrics: ['Financed Emissions', 'Portfolio Analysis', 'TCFD Disclosure'], icon: ShieldIcon },
              { title: 'Technology', metrics: ['Data Center PUE', 'Renewable Procurement', 'Scope 3 Category 1'], icon: GlobeIcon },
            ].map((solution, i) => (
              <Card key={i} className="h-full">
                <CardBody className="p-6">
                  <div className="p-3 bg-primary-100 dark:bg-primary-900/30 rounded-xl w-fit mb-4">
                    <solution.icon className="h-6 w-6 text-primary-600 dark:text-primary-400" />
                  </div>
                  <h3 className="text-xl font-semibold text-secondary-900 dark:text-white mb-3">{solution.title}</h3>
                  <ul className="space-y-2">
                    {solution.metrics.map((metric, j) => (
                      <li key={j} className="flex items-center gap-2 text-sm text-secondary-600 dark:text-secondary-400">
                        <CheckIcon className="h-4 w-4 text-green-500" />
                        {metric}
                      </li>
                    ))}
                  </ul>
                </CardBody>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* Integrations Section */}
      <section id="integrations" className="py-20 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto">
          <div className="text-center mb-16">
            <Badge variant="outline" className="mb-4">Integrations</Badge>
            <h2 className="text-3xl sm:text-4xl font-bold text-secondary-900 dark:text-white mb-4">
              Connect Your Entire Data Ecosystem
            </h2>
            <p className="text-lg text-secondary-600 dark:text-secondary-300 max-w-2xl mx-auto">
              Pre-built connectors for seamless data ingestion from 300+ sources.
            </p>
          </div>
          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-4">
            {integrations.map((integration, i) => (
              <Card key={i} className="text-center p-6 hover:border-primary-300 dark:hover:border-primary-700 transition-colors">
                <p className="text-3xl font-bold text-secondary-900 dark:text-white">{integration.count}+</p>
                <p className="text-secondary-600 dark:text-secondary-400 mt-1">{integration.name}</p>
              </Card>
            ))}
          </div>
          <div className="text-center mt-10">
            <Button variant="outline" size="lg" asChild>
              <Link to="/integrations">View All Integrations <ChevronRightIcon className="h-4 w-4 ml-2" /></Link>
            </Button>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 px-4 sm:px-6 lg:px-8 bg-secondary-900 dark:bg-secondary-950">
        <div className="max-w-4xl mx-auto text-center">
          <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
            Ready to Transform Your Climate Strategy?
          </h2>
          <p className="text-lg text-secondary-300 mb-8">
            Join 500+ organizations using Green Combats to achieve their sustainability goals.
          </p>
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link to="/auth/register">
              <Button variant="primary" size="lg" className="w-full sm:w-auto gap-2 bg-primary-600 hover:bg-primary-700">
                Start Free Trial
                <ArrowRightIcon className="h-5 w-5" />
              </Button>
            </Link>
            <Link to="/contact">
              <Button variant="outline" size="lg" className="w-full sm:w-auto border-secondary-600 text-secondary-200 hover:bg-secondary-800">
                Contact Sales
              </Button>
            </Link>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-secondary-950 dark:bg-black border-t border-secondary-800 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8 mb-12">
            <div className="col-span-2 md:col-span-1">
              <Link to="/" className="flex items-center gap-2 mb-4">
                <div className="h-10 w-10 rounded-xl bg-primary-600 flex items-center justify-center">
                  <span className="text-white font-bold text-xl">GC</span>
                </div>
                <span className="text-xl font-bold text-white">Green Combats</span>
              </Link>
              <p className="text-secondary-400 text-sm">Climate action intelligence platform for a sustainable future.</p>
            </div>
            <div>
              <h4 className="font-semibold text-white mb-4">Product</h4>
              <ul className="space-y-2 text-sm text-secondary-400">
                <li><Link to="/features" className="hover:text-white transition-colors">Features</Link></li>
                <li><Link to="/pricing" className="hover:text-white transition-colors">Pricing</Link></li>
                <li><Link to="/integrations" className="hover:text-white transition-colors">Integrations</Link></li>
                <li><Link to="/changelog" className="hover:text-white transition-colors">Changelog</Link></li>
                <li><Link to="/roadmap" className="hover:text-white transition-colors">Roadmap</Link></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold text-white mb-4">Resources</h4>
              <ul className="space-y-2 text-sm text-secondary-400">
                <li><Link to="/docs" className="hover:text-white transition-colors">Documentation</Link></li>
                <li><Link to="/blog" className="hover:text-white transition-colors">Blog</Link></li>
                <li><Link to="/webinars" className="hover:text-white transition-colors">Webinars</Link></li>
                <li><Link to="/case-studies" className="hover:text-white transition-colors">Case Studies</Link></li>
                <li><Link to="/help" className="hover:text-white transition-colors">Help Center</Link></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold text-white mb-4">Company</h4>
              <ul className="space-y-2 text-sm text-secondary-400">
                <li><Link to="/about" className="hover:text-white transition-colors">About</Link></li>
                <li><Link to="/careers" className="hover:text-white transition-colors">Careers</Link></li>
                <li><Link to="/press" className="hover:text-white transition-colors">Press</Link></li>
                <li><Link to="/contact" className="hover:text-white transition-colors">Contact</Link></li>
                <li><Link to="/partners" className="hover:text-white transition-colors">Partners</Link></li>
              </ul>
            </div>
          </div>
          <div className="pt-8 border-t border-secondary-800 flex flex-col md:flex-row items-center justify-between gap-4">
            <p className="text-sm text-secondary-500">© 2024 Green Combats. All rights reserved.</p>
            <div className="flex items-center gap-6 text-sm text-secondary-500">
              <Link to="/privacy" className="hover:text-white transition-colors">Privacy</Link>
              <Link to="/terms" className="hover:text-white transition-colors">Terms</Link>
              <Link to="/cookies" className="hover:text-white transition-colors">Cookies</Link>
              <Link to="/security" className="hover:text-white transition-colors">Security</Link>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}