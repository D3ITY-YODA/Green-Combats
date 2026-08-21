import { Link } from 'react-router-dom';
import { HomeIcon, SearchIcon, ArrowLeftIcon, CompassIcon, FileQuestionIcon } from 'lucide-react';
import { Button } from '../../components/ui/Button';

export default function NotFoundPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-secondary-50 dark:bg-secondary-900 px-4">
      <div className="text-center max-w-md">
        <div className="mb-8">
          <span className="text-9xl font-bold text-primary-600/20 dark:text-primary-400/20">404</span>
        </div>
        <h1 className="text-3xl font-bold text-secondary-900 dark:text-white mb-4">Page Not Found</h1>
        <p className="text-lg text-secondary-600 dark:text-secondary-400 mb-8">
          Sorry, we couldn't find the page you're looking for. It might have been moved or doesn't exist.
        </p>
        <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
          <Button variant="primary" asChild>
            <Link to="/">
              <HomeIcon className="h-4 w-4 mr-2" />
              Go Home
            </Link>
          </Button>
          <Button variant="outline" asChild>
            <Link to="/dashboard">
              <CompassIcon className="h-4 w-4 mr-2" />
              Dashboard
            </Link>
          </Button>
        </div>
        <div className="mt-12 space-y-4 text-sm text-secondary-500 dark:text-secondary-400">
          <p>Or try one of these:</p>
          <div className="flex flex-wrap items-center justify-center gap-4">
            <Link to="/help" className="flex items-center gap-1.5 hover:text-primary-600 dark:hover:text-primary-400 transition-colors">
              <FileQuestionIcon className="h-4 w-4" />
              Help Center
            </Link>
            <span className="text-secondary-300 dark:text-secondary-600">·</span>
            <Link to="/docs" className="hover:text-primary-600 dark:hover:text-primary-400 transition-colors">Documentation</Link>
            <span className="text-secondary-300 dark:text-secondary-600">·</span>
            <Link to="/contact" className="hover:text-primary-600 dark:hover:text-primary-400 transition-colors">Contact Support</Link>
          </div>
        </div>
      </div>
    </div>
  );
}