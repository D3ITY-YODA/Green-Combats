import { forwardRef, InputHTMLAttributes, LabelHTMLAttributes, TextareaHTMLAttributes, SelectHTMLAttributes } from 'react';
import { cn } from '@utils/helpers';

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  hint?: string;
  iconLeft?: React.ReactNode;
  iconRight?: React.ReactNode;
  fullWidth?: boolean;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, hint, iconLeft, iconRight, fullWidth = true, className, id, ...props }, ref) => {
    const inputId = id || `input-${Math.random().toString(36).substr(2, 9)}`;
    const hintId = hint ? `${inputId}-hint` : undefined;
    const errorId = error ? `${inputId}-error` : undefined;
    
    return (
      <div className={cn(fullWidth ? 'w-full' : '', className)}>
        {label && (
          <label htmlFor={inputId} className="block text-sm font-medium text-secondary-700 dark:text-secondary-300 mb-1.5">
            {label}
          </label>
        )}
        <div className="relative">
          {iconLeft && (
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-secondary-400">
              {iconLeft}
            </div>
          )}
          <input
            ref={ref}
            id={inputId}
            className={cn(
              'w-full rounded-lg border bg-white px-4 py-2.5 text-sm placeholder:text-secondary-400 transition-colors',
              'focus:outline-none focus:ring-2 focus:ring-primary-500/20',
              iconLeft ? 'pl-10' : '',
              iconRight ? 'pr-10' : '',
              error
                ? 'border-red-500 focus:border-red-500 focus:ring-red-500/20'
                : 'border-secondary-300 focus:border-primary-500 dark:border-secondary-600 dark:bg-secondary-900 dark:text-secondary-100',
              props.disabled && 'opacity-50 cursor-not-allowed bg-secondary-50 dark:bg-secondary-800'
            )}
            aria-describedby={cn(hintId, errorId)}
            aria-invalid={error ? 'true' : 'false'}
            {...props}
          />
          {iconRight && (
            <div className="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none text-secondary-400">
              {iconRight}
            </div>
          )}
        </div>
        {error && (
          <p id={errorId} className="mt-1.5 text-sm text-red-600 dark:text-red-400" role="alert">
            {error}
          </p>
        )}
        {hint && !error && (
          <p id={hintId} className="mt-1.5 text-sm text-secondary-500 dark:text-secondary-400">
            {hint}
          </p>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';

export interface TextareaProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string;
  error?: string;
  hint?: string;
  fullWidth?: boolean;
}

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ label, error, hint, fullWidth = true, className, id, ...props }, ref) => {
    const textareaId = id || `textarea-${Math.random().toString(36).substr(2, 9)}`;
    const hintId = hint ? `${textareaId}-hint` : undefined;
    const errorId = error ? `${textareaId}-error` : undefined;
    
    return (
      <div className={cn(fullWidth ? 'w-full' : '', className)}>
        {label && (
          <label htmlFor={textareaId} className="block text-sm font-medium text-secondary-700 dark:text-secondary-300 mb-1.5">
            {label}
          </label>
        )}
        <textarea
          ref={ref}
          id={textareaId}
          className={cn(
            'w-full rounded-lg border bg-white px-4 py-2.5 text-sm placeholder:text-secondary-400 transition-colors resize-y min-h-[100px]',
            'focus:outline-none focus:ring-2 focus:ring-primary-500/20 focus:border-primary-500',
            error
              ? 'border-red-500 focus:border-red-500 focus:ring-red-500/20'
              : 'border-secondary-300 dark:border-secondary-600 dark:bg-secondary-900 dark:text-secondary-100',
            props.disabled && 'opacity-50 cursor-not-allowed bg-secondary-50 dark:bg-secondary-800'
          )}
          aria-describedby={cn(hintId, errorId)}
          aria-invalid={error ? 'true' : 'false'}
          {...props}
        />
        {error && (
          <p id={errorId} className="mt-1.5 text-sm text-red-600 dark:text-red-400" role="alert">
            {error}
          </p>
        )}
        {hint && !error && (
          <p id={hintId} className="mt-1.5 text-sm text-secondary-500 dark:text-secondary-400">
            {hint}
          </p>
        )}
      </div>
    );
  }
);

Textarea.displayName = 'Textarea';

export interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
  label?: string;
  error?: string;
  hint?: string;
  options: { value: string; label: string; disabled?: boolean }[];
  placeholder?: string;
  fullWidth?: boolean;
}

export const Select = forwardRef<HTMLSelectElement, SelectProps>(
  ({ label, error, hint, options, placeholder, fullWidth = true, className, id, ...props }, ref) => {
    const selectId = id || `select-${Math.random().toString(36).substr(2, 9)}`;
    const hintId = hint ? `${selectId}-hint` : undefined;
    const errorId = error ? `${selectId}-error` : undefined;
    
    return (
      <div className={cn(fullWidth ? 'w-full' : '', className)}>
        {label && (
          <label htmlFor={selectId} className="block text-sm font-medium text-secondary-700 dark:text-secondary-300 mb-1.5">
            {label}
          </label>
        )}
        <select
          ref={ref}
          id={selectId}
          className={cn(
            'w-full rounded-lg border bg-white px-4 py-2.5 text-sm transition-colors appearance-none',
            'focus:outline-none focus:ring-2 focus:ring-primary-500/20 focus:border-primary-500',
            error
              ? 'border-red-500 focus:border-red-500 focus:ring-red-500/20'
              : 'border-secondary-300 dark:border-secondary-600 dark:bg-secondary-900 dark:text-secondary-100',
            props.disabled && 'opacity-50 cursor-not-allowed bg-secondary-50 dark:bg-secondary-800'
          )}
          aria-describedby={cn(hintId, errorId)}
          aria-invalid={error ? 'true' : 'false'}
          {...props}
        >
          {placeholder && (
            <option value="" disabled>
              {placeholder}
            </option>
          )}
          {options.map((option) => (
            <option key={option.value} value={option.value} disabled={option.disabled}>
              {option.label}
            </option>
          ))}
        </select>
        {error && (
          <p id={errorId} className="mt-1.5 text-sm text-red-600 dark:text-red-400" role="alert">
            {error}
          </p>
        )}
        {hint && !error && (
          <p id={hintId} className="mt-1.5 text-sm text-secondary-500 dark:text-secondary-400">
            {hint}
          </p>
        )}
      </div>
    );
  }
);

Select.displayName = 'Select';

export interface LabelProps extends LabelHTMLAttributes<HTMLLabelElement> {
  required?: boolean;
}

export const Label = forwardRef<HTMLLabelElement, LabelProps>(
  ({ children, required, className, ...props }, ref) => (
    <label
      ref={ref}
      className={cn('block text-sm font-medium text-secondary-700 dark:text-secondary-300 mb-1.5', className)}
      {...props}
    >
      {children}
      {required && <span className="text-red-500 ml-1" aria-hidden="true">*</span>}
    </label>
  )
);

Label.displayName = 'Label';

export interface CheckboxProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  description?: string;
}

export const Checkbox = forwardRef<HTMLInputElement, CheckboxProps>(
  ({ label, description, className, id, ...props }, ref) => {
    const checkboxId = id || `checkbox-${Math.random().toString(36).substr(2, 9)}`;
    
    return (
      <div className="flex items-start gap-3">
        <input
          ref={ref}
          type="checkbox"
          id={checkboxId}
          className={cn(
            'mt-0.5 h-4 w-4 rounded border-secondary-300 text-primary-600',
            'focus:ring-2 focus:ring-primary-500/20 focus:ring-offset-2',
            'disabled:opacity-50 disabled:cursor-not-allowed',
            'dark:border-secondary-600 dark:bg-secondary-800',
            className
          )}
          {...props}
        />
        {(label || description) && (
          <div className="flex flex-col">
            {label && (
              <label htmlFor={checkboxId} className="text-sm text-secondary-900 dark:text-secondary-100 cursor-pointer">
                {label}
              </label>
            )}
            {description && (
              <p className="text-sm text-secondary-500 dark:text-secondary-400">{description}</p>
            )}
          </div>
        )}
      </div>
    );
  }
);

Checkbox.displayName = 'Checkbox';

export interface RadioProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  description?: string;
}

export const Radio = forwardRef<HTMLInputElement, RadioProps>(
  ({ label, description, className, id, ...props }, ref) => {
    const radioId = id || `radio-${Math.random().toString(36).substr(2, 9)}`;
    
    return (
      <div className="flex items-start gap-3">
        <input
          ref={ref}
          type="radio"
          id={radioId}
          className={cn(
            'mt-0.5 h-4 w-4 border-secondary-300 text-primary-600',
            'focus:ring-2 focus:ring-primary-500/20 focus:ring-offset-2',
            'disabled:opacity-50 disabled:cursor-not-allowed',
            'dark:border-secondary-600 dark:bg-secondary-800',
            className
          )}
          {...props}
        />
        {(label || description) && (
          <div className="flex flex-col">
            {label && (
              <label htmlFor={radioId} className="text-sm text-secondary-900 dark:text-secondary-100 cursor-pointer">
                {label}
              </label>
            )}
            {description && (
              <p className="text-sm text-secondary-500 dark:text-secondary-400">{description}</p>
            )}
          </div>
        )}
      </div>
    );
  }
);

Radio.displayName = 'Radio';

export interface SwitchProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  description?: string;
}

export const Switch = forwardRef<HTMLInputElement, SwitchProps>(
  ({ label, description, className, id, ...props }, ref) => {
    const switchId = id || `switch-${Math.random().toString(36).substr(2, 9)}`;
    
    return (
      <div className="flex items-center gap-3">
        <label htmlFor={switchId} className="relative inline-flex items-center cursor-pointer">
          <input
            ref={ref}
            type="checkbox"
            id={switchId}
            className={cn(
              'sr-only peer',
              'focus-visible:ring-2 focus-visible:ring-primary-500/20 focus-visible:ring-offset-2',
              className
            )}
            {...props}
          />
          <div className={cn(
            'w-11 h-6 rounded-full border-2 transition-colors duration-200',
            'peer-focus-visible:ring-2 peer-focus-visible:ring-primary-500/20 peer-focus-visible:ring-offset-2',
            'peer-checked:bg-primary-600 peer-checked:border-primary-600',
            'peer-checked:after:translate-x-full',
            'peer-disabled:opacity-50 peer-disabled:cursor-not-allowed',
            props.checked ? 'bg-primary-600 border-primary-600' : 'bg-white border-secondary-300 dark:border-secondary-600'
          )}>
            <div className={cn(
              'w-5 h-5 rounded-full bg-white shadow-lg transform transition-transform duration-200',
              'peer-checked:translate-x-full'
            )} />
          </div>
        </label>
        {(label || description) && (
          <div className="flex flex-col">
            {label && <span className="text-sm text-secondary-900 dark:text-secondary-100">{label}</span>}
            {description && <p className="text-sm text-secondary-500 dark:text-secondary-400">{description}</p>}
          </div>
        )}
      </div>
    );
  }
);

Switch.displayName = 'Switch';