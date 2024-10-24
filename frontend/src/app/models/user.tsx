import mongoose, { Schema, Document, Model } from 'mongoose';

export interface IUser extends Document {
    _id: string;
  name: string;
  email: string;
  password: string;
  role: string;
}

const UserSchema: Schema = new mongoose.Schema({
  name: { type: String, required: true },
  email: { type: String, required: true, unique: true },
  password: { type: String, required: true },
  role: { type: String, default: 'user' },
});

const User: Model<IUser> = mongoose.models.User || mongoose.model<IUser>('User', UserSchema);
export default User;
