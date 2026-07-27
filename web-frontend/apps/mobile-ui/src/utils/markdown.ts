import DOMPurify from 'dompurify';
import { marked } from 'marked';

marked.setOptions({ gfm: true, breaks: true });

export function renderMarkdown(md: string): string {
  if (!md) return '';
  try {
    const raw = marked.parse(md) as string;
    return DOMPurify.sanitize(raw);
  } catch {
    return String(md);
  }
}

export function escapeHtml(text: string): string {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}
