import { NavigationBar } from '@/features/NavigationBar';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { basename } from './const';
import { DetailPage } from './pages/DetailPage';
import { HomePage } from './pages/Home.page';

const router = createBrowserRouter(
  [
    {
      path: '/',
      element: (
        <NavigationBar>
          <HomePage />
        </NavigationBar>
      ),
    },
    {
      path: '/detail',
      element: (
        <NavigationBar>
          <DetailPage />
        </NavigationBar>
      ),
    },
  ],
  { basename }
);

export function Router() {
  return <RouterProvider router={router} />;
}
