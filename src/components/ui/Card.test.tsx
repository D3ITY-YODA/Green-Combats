import { render, screen, cleanup } from '@testing-library/react'
import { describe, it, expect, afterEach } from 'vitest'
import { Card, CardHeader, CardBody, CardFooter, CardTitle, CardDescription } from './Card'

describe('Card', () => {
  afterEach(() => {
    cleanup()
  })

  it('renders card with children', () => {
    render(<Card>Card content</Card>)
    expect(screen.getByText(/card content/i)).toBeInTheDocument()
  })

  it('applies default variant classes', () => {
    const { container } = render(<Card>Content</Card>)
    const card = container.firstChild as HTMLElement
    expect(card).toHaveClass('border')
    expect(card).toHaveClass('border-secondary-200')
    expect(card).toHaveClass('bg-white')
    expect(card).toHaveClass('shadow-sm')
    expect(card).toHaveClass('rounded-xl')
  })

  it('applies outlined variant classes', () => {
    const { container } = render(<Card variant="outlined">Content</Card>)
    const card = container.firstChild as HTMLElement
    expect(card).toHaveClass('border-2')
  })

  it('applies elevated variant classes', () => {
    const { container } = render(<Card variant="elevated">Content</Card>)
    const card = container.firstChild as HTMLElement
    expect(card).toHaveClass('border-none')
    expect(card).toHaveClass('shadow-lg')
  })

  it('applies padding classes', () => {
    const { container, unmount } = render(<Card padding="none">Content</Card>)
    const card = container.firstChild as HTMLElement
    expect(card).not.toHaveClass('p-4')
    expect(card).not.toHaveClass('p-6')
    expect(card).not.toHaveClass('p-8')
    unmount()

    const { container: container2, unmount: unmount2 } = render(<Card padding="sm">Content</Card>)
    expect(container2.firstChild).toHaveClass('p-4')
    unmount2()

    const { container: container3, unmount: unmount3 } = render(<Card padding="md">Content</Card>)
    expect(container3.firstChild).toHaveClass('p-6')
    unmount3()

    const { container: container4, unmount: unmount4 } = render(<Card padding="lg">Content</Card>)
    expect(container4.firstChild).toHaveClass('p-8')
    unmount4()
  })

  it('applies hover effect when hover prop is true', () => {
    const { container } = render(<Card hover>Content</Card>)
    const card = container.firstChild as HTMLElement
    expect(card).toHaveClass('transition-shadow')
    expect(card).toHaveClass('hover:shadow-md')
    expect(card).toHaveClass('cursor-pointer')
  })

  it('renders CardHeader with border', () => {
    const { container } = render(
      <Card>
        <CardHeader>Header content</CardHeader>
      </Card>
    )
    expect(screen.getByText(/header content/i)).toBeInTheDocument()
    const header = container.querySelector('[class*="border-b"]')
    expect(header).toHaveClass('border-b')
    expect(header).toHaveClass('px-6')
    expect(header).toHaveClass('py-4')
  })

  it('renders CardBody with padding', () => {
    const { container } = render(
      <Card>
        <CardBody>Body content</CardBody>
      </Card>
    )
    expect(screen.getByText(/body content/i)).toBeInTheDocument()
    const body = container.querySelector('[class*="p-6"]')
    expect(body).toHaveClass('p-6')
  })

  it('renders CardFooter with border', () => {
    const { container } = render(
      <Card>
        <CardFooter>Footer content</CardFooter>
      </Card>
    )
    expect(screen.getByText(/footer content/i)).toBeInTheDocument()
    const footer = container.querySelector('[class*="border-t"]')
    expect(footer).toHaveClass('border-t')
    expect(footer).toHaveClass('px-6')
    expect(footer).toHaveClass('py-4')
  })

  it('renders CardTitle with correct styling', () => {
    render(
      <Card>
        <CardTitle>Card Title</CardTitle>
      </Card>
    )
    expect(screen.getByText(/card title/i)).toBeInTheDocument()
    const title = screen.getByText(/card title/i)
    expect(title).toHaveClass('text-lg')
    expect(title).toHaveClass('font-semibold')
    expect(title).toHaveClass('text-secondary-900')
  })

  it('renders CardTitle with custom as prop', () => {
    render(
      <Card>
        <CardTitle as="h2">H2 Title</CardTitle>
      </Card>
    )
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent(/h2 title/i)
  })

  it('renders CardDescription with correct styling', () => {
    render(
      <Card>
        <CardDescription>Card description</CardDescription>
      </Card>
    )
    expect(screen.getByText(/card description/i)).toBeInTheDocument()
    const desc = screen.getByText(/card description/i)
    expect(desc).toHaveClass('text-sm')
    expect(desc).toHaveClass('text-secondary-500')
    expect(desc).toHaveClass('mt-1')
  })
})