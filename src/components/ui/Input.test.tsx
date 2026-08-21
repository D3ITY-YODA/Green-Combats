import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import { Input } from './Input'

describe('Input', () => {
  it('renders input with label', () => {
    render(<Input label="Email" placeholder="Enter email" />)
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument()
    expect(screen.getByPlaceholderText(/enter email/i)).toBeInTheDocument()
  })

  it('renders error message when error prop is provided', () => {
    render(<Input label="Email" error="Invalid email" />)
    expect(screen.getByRole('alert')).toHaveTextContent(/invalid email/i)
  })

  it('applies error styling when error prop is provided', () => {
    render(<Input label="Email" error="Invalid email" />)
    const input = screen.getByLabelText(/email/i)
    expect(input).toHaveClass('border-red-500')
    expect(input).toHaveAttribute('aria-invalid', 'true')
  })

  it('renders hint text when hint prop is provided', () => {
    render(<Input label="Password" hint="Must be 8+ characters" />)
    expect(screen.getByText(/must be 8\+ characters/i)).toBeInTheDocument()
  })

  it('does not show hint when error is present', () => {
    render(<Input label="Email" hint="Enter your email" error="Invalid" />)
    expect(screen.queryByText(/enter your email/i)).not.toBeInTheDocument()
  })

  it('calls onChange handler', () => {
    const handleChange = vi.fn()
    render(<Input onChange={handleChange} />)
    fireEvent.change(screen.getByRole('textbox'), { target: { value: 'test' } })
    expect(handleChange).toHaveBeenCalled()
  })

  it('disables input when disabled prop is true', () => {
    render(<Input disabled />)
    expect(screen.getByRole('textbox')).toBeDisabled()
  })

  it('applies iconLeft styling', () => {
    const TestIcon = () => <svg data-testid="left-icon" />
    render(<Input iconLeft={<TestIcon />} />)
    expect(screen.getByTestId('left-icon')).toBeInTheDocument()
    expect(screen.getByRole('textbox')).toHaveClass('pl-10')
  })

  it('applies iconRight styling', () => {
    const TestIcon = () => <svg data-testid="right-icon" />
    render(<Input iconRight={<TestIcon />} />)
    expect(screen.getByTestId('right-icon')).toBeInTheDocument()
    expect(screen.getByRole('textbox')).toHaveClass('pr-10')
  })

  it('renders with correct id when provided', () => {
    render(<Input id="custom-id" />)
    expect(screen.getByRole('textbox')).toHaveAttribute('id', 'custom-id')
  })

  it('generates id when not provided', () => {
    render(<Input label="Email" />)
    const input = screen.getByLabelText(/email/i)
    expect(input).toHaveAttribute('id')
    expect(input.id).toMatch(/^input-/)
  })
})